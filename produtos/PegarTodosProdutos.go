package produtos

import (
	"log"
	"net/http"
	interno "sos/backend/interno/db"
	"strings"

	"firebase.google.com/go/v4/auth"
)

func PegarTodosProdutos(db *interno.Queries, app *auth.Client) http.HandlerFunc {
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
		dd := user.Claims["role"]

		log.Println(dd)

		w.WriteHeader(200)
		w.Write([]byte(`{"message": "Hello World"}`))
	}
}
