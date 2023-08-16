package auth

import (
	"net/http"
	interno "sos/backend/interno/db"

	"firebase.google.com/go/v4/auth"
)

func Login(db *interno.Queries, app *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain")
		_, _, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized"))
			return
		}

		w.Write([]byte("login!"))
	}
}
