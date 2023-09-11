package produtos

import (
	"encoding/json"
	"log"
	"net/http"
	interno "sos/backend/interno/db"
	"strings"

	"cloud.google.com/go/storage"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func PegarProdutoUnico(db *interno.Queries, app *auth.Client, bucket *storage.BucketHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tokenC, err := r.Cookie("session")

		token := r.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			if err != nil {
				log.Println(err)
				w.WriteHeader(401)
				w.Write([]byte(`{"message": "Token não fornecido"}`))
				return

			}
			token = tokenC.Value
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

		id := chi.URLParam(r, "id")
		if id == "" {
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "ID não fornecido"}`))
			return
		}

		produto, err := db.PegarProdutoUnico(r.Context(), uuid.MustParse(id))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "Produto não encontrado"}`))
			return
		}

		json.NewEncoder(w).Encode(produto)
		w.WriteHeader(201)
	}
}
