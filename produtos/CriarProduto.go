package produtos

import (
	"log"
	"net/http"
	interno "sos/backend/interno/db"
	"strings"

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
	Imagem           string         `json:"imagem"`
	Categoria        []string       `json:"categoria"`
}

func CriaProduto(db *interno.Queries, app *auth.Client) http.HandlerFunc {
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

		idProduto, err := db.CriarProduto(r.Context(), interno.CriarProdutoParams{
			Nome:             produto.Nome,
			Descricao:        produto.Descricao,
			Preco:            produto.Preco,
			UnidadeMedida:    interno.TipoUnidadeMedida(produto.UnidadeMedida),
			QuantidadePacote: produto.QuantidadePacote,
			Peso:             produto.Peso,
			Ativo:            produto.Ativo,
			Ordem:            produto.Ordem,
			Imagem:           produto.Imagem,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao criar produto"}`))
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
				inserirCategoriaNoProdutoParams = append(inserirCategoriaNoProdutoParams, interno.InserirCategoriaNoProdutoParams{A: idProduto, B: idCategoria})
			}
		})

		categoriaInserida := db.InserirCategoriaNoProduto(r.Context(), inserirCategoriaNoProdutoParams)
		defer categoriaInserida.Close()

		w.WriteHeader(200)
	}
}
