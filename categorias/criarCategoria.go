package categorias

import (
	"encoding/json"
	"log"
	"net/http"
	interno "sos/backend/interno/db"
	"sos/backend/model"
	"strings"

	"firebase.google.com/go/v4/auth"
)

type categoria struct {
	Nome   string `json:"nome"`
	Imagem string `json:"imagem"`
}

func CriarCategoria(db *interno.Queries, app *auth.Client) http.HandlerFunc {
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

		var categoria categoria
		json.NewDecoder(r.Body).Decode(&categoria)

		if categoria.Nome == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Nome não pode ser vazio"})
			return
		}

		if categoria.Imagem == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Imagem não pode ser vazio"})
			return
		}

		_, err = db.CriarCategoria(r.Context(), interno.CriarCategoriaParams{Nome: categoria.Nome, Imagem: categoria.Imagem})
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 500, Mensagem: "Erro ao criar categoria"})
			return
		}

		w.WriteHeader(201)
	}
}
