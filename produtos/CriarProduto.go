package produtos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	interno "sos/backend/interno/db"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"cloud.google.com/go/storage"
	"firebase.google.com/go/v4/auth"
	"github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/image/draw"
)

type Produto struct {
	Nome             string         `json:"nome"`
	Descricao        string         `json:"descricao"`
	Preco            pgtype.Numeric `json:"preco"`
	UnidadeMedida    string         `json:"unidadeMedida"`
	QuantidadePacote int32          `json:"quantidadePacote"`
	Peso             pgtype.Numeric `json:"peso"`
	Ativo            bool           `json:"ativo"`
	Ordem            int32          `json:"ordem"`
	Imagem           http.File      `json:"imagem"`
	Categoria        []string       `json:"categoria"`
}

func CriaProduto(db *interno.Queries, app *auth.Client, bucket *storage.BucketHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)

		if token == "" {
			w.WriteHeader(401)
			w.Write([]byte(`{"message": "Token não fornecido"}`))
			return

		}

		user, err := app.VerifyIDToken(r.Context(), token)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(`{"message": "Token inválido"}`))
			return
		}

		role := user.Claims["role"]
		log.Println(role)
		// if role != "ADMIN" && role != "DONO" {
		// 	w.WriteHeader(401)
		// 	w.Write([]byte(`{"message": "Usuário não autorizado"}`))
		// 	return
		// }

		var produto Produto
		jsonData := r.FormValue("produto")
		json.Unmarshal([]byte(jsonData), &produto)

		idProduto, err := db.CriarProduto(r.Context(), interno.CriarProdutoParams{
			Nome:             produto.Nome,
			Descricao:        produto.Descricao,
			Preco:            produto.Preco,
			UnidadeMedida:    interno.TipoUnidadeMedida(produto.UnidadeMedida),
			QuantidadePacote: produto.QuantidadePacote,
			Peso:             produto.Peso,
			Ativo:            produto.Ativo,
			Ordem:            produto.Ordem,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao criar produto"}`))
			return
		}
		imagem, _, err := r.FormFile("imagem")
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao pegar imagem"}`))
			return
		}
		defer imagem.Close()
		nomeImagemProduto := "produtos/" + idProduto.String() + ".webp"

		imgBytes, err := converterParaWebp(imagem) //io.ReadAll(imagem)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao ler imagem"}`))
			return
		}
		obj := bucket.Object(nomeImagemProduto)
		write := obj.NewWriter(r.Context())
		defer write.Close()

		_, err = io.Copy(write, bytes.NewReader(imgBytes))
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao salvar imagem"}`))
			return
		}

		err = db.AtualizarImagemProduto(r.Context(), idProduto)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao atualizar imagem do produto"}`))
			return
		}

		var inserirCategoriaNoProdutoParams []interno.InserirCategoriaNoProdutoParams

		categoria := db.PegarCategoriaPeloNome(r.Context(), produto.Categoria)
		defer categoria.Close()
		categoria.Query(func(i int, id []uuid.UUID, err error) {
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				w.Write([]byte(`{"message": "Erro ao pegar categoria"}`))
				return
			}
			for _, idCategoria := range id {
				inserirCategoriaNoProdutoParams = append(inserirCategoriaNoProdutoParams, interno.InserirCategoriaNoProdutoParams{A: idCategoria, B: idProduto})
			}
		})

		categoriaInserida := db.InserirCategoriaNoProduto(r.Context(), inserirCategoriaNoProdutoParams)
		defer categoriaInserida.Close()
		categoriaInserida.Exec(func(i int, err error) {
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				w.Write([]byte(`{"message": "Erro ao inserir categoria no produto"}`))
				return
			}
		})

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(""))
	}
}

func converterParaWebp(file multipart.File) (webpBytes []byte, err error) {
	// Decodificar a imagem no formato apropriado (JPEG, PNG, GIF, etc.)
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar o início do arquivo: %v", err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar a imagem: %v", err)
	}
	defer file.Close()
	img = resizeImage(img, 512, 512)
	// Converter a imagem para WebP otimizado
	var webpBuffer bytes.Buffer
	err = webp.Encode(&webpBuffer, img, &webp.Options{Quality: 80, Lossless: false})
	if err != nil {
		return nil, fmt.Errorf("erro ao converter a imagem para WebP: %v", err)
	}

	// Armazenar o resultado como []byte em uma variável
	webpBytes = webpBuffer.Bytes()

	return webpBytes, nil
}

func resizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(newImage, newImage.Bounds(), img, bounds, draw.Src, nil)
	return newImage
}
