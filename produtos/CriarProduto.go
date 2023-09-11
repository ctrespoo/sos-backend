package produtos

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	interno "sos/backend/interno/db"
	"strings"

	"cloud.google.com/go/storage"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		nomeImagemProduto := "produtos/" + idProduto.String() + ".png"
		imgBytes, err := io.ReadAll(imagem)
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

		w.Write([]byte(""))
		w.WriteHeader(201)
	}
}
