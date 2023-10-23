package produtos

import (
	"encoding/json"
	"net/http"
	interno "sos/backend/interno/db"

	"firebase.google.com/go/v4/auth"
)

func PegarTodosProdutos(db *interno.Queries, app *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// token := r.Header.Get("Authorization")
		// token = strings.Replace(token, "Bearer ", "", 1)
		// if token == "" {
		// 	w.WriteHeader(401)
		// 	w.Write([]byte(`{"message": "Token não fornecido"}`))
		// 	return
		// }

		// user, err := app.VerifyIDToken(r.Context(), token)
		// if err != nil {
		// 	w.WriteHeader(401)
		// 	w.Write([]byte(`{"message": "Token inválido"}`))
		// 	return
		// }

		// role := user.Claims["role"]
		// log.Println(role)
		// if role != "ADMIN" && role != "DONO" {
		// 	w.WriteHeader(401)
		// 	w.Write([]byte(`{"message": "Usuário não autorizado"}`))
		// 	return
		// }

		produtos, err := db.PegarTodosProdutos(r.Context(), interno.PegarTodosProdutosParams{Limit: 10, Offset: 0})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao pegar produtos"}`))
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(produtos)
	}
}
