package produtos

import (
	"encoding/json"
	"log"
	"net/http"
	interno "sos/backend/interno/db"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func AtualizarProduto(db *interno.Queries, app *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		tokenC, err := r.Cookie("session")
		if err != nil {
			log.Println(err)
			w.WriteHeader(401)
			w.Write([]byte(`{"message": "Token não fornecido"}`))
			return
		}

		token := r.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			token = tokenC.Value

			if token == "" {
				w.WriteHeader(401)
				w.Write([]byte(`{"message": "Token não fornecido"}`))
				return

			}
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

		var nomesCategorias []string
		var categoriasParaExcluir []string

		json.NewDecoder(r.Body).Decode(&nomesCategorias)

		id := chi.URLParam(r, "id")

		categorias, err := db.PegarCategoriasProduto(r.Context(), uuid.MustParse(id))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "Produto não encontrado"}`))
			return
		}

		verificarCategorias(&categorias, &nomesCategorias, &categoriasParaExcluir)
		log.Println(categoriasParaExcluir)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(categoriasParaExcluir)

	}
}

func verificarCategorias(categorias *[]interno.Categoria, nomesCategorias *[]string, categoriasParaExcluir *[]string) {
	// Percorra as categorias
	for _, categoria := range *categorias {
		// Verifique se o nome da categoria não está na lista de nomesCategorias
		if !contemNomeCategoria(categoria.Nome, nomesCategorias) {
			// Adicione o nome da categoria à lista de categoriasParaExcluir
			*categoriasParaExcluir = append(*categoriasParaExcluir, categoria.Nome)
		}
	}
}

func contemNomeCategoria(nome string, nomesCategorias *[]string) bool {
	// Percorra a lista de nomesCategorias e verifique se o nome está presente
	for _, nomeCategoria := range *nomesCategorias {
		if nomeCategoria == nome {
			return true
		}
	}
	return false
}
