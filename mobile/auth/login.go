package mobile

import (
	"net/http"
	interno "sos/backend/interno/db"
)

func Login(db *interno.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("login!"))
	}
}
