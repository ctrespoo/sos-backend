package mobile

import (
	"net/http"
	interno "sos/backend/interno/db"
)

func Login(db *interno.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain")
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized"))
			return
		}

		w.Write([]byte("login!"))
	}
}
