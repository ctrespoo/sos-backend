package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	auth "sos/backend/mobile/auth"

	interno "sos/backend/interno/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {

	dbs := os.Getenv("DATABASE_URL")

	if dbs == "" {
		log.Fatal("DATABASE_URL não definida")
	}

	db, err := sql.Open("postgres", dbs)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dbtx := interno.New(db)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/mobile/api/v1", func(r chi.Router) {
		r.Post("/auth/login", auth.Login(dbtx))
		r.Post("/auth/registrar", auth.CriarConta(dbtx))
		r.Group(func(r chi.Router) {

		})
	})

	srv := &http.Server{
		Addr: "0.0.0.0:20000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Println("Listening on 20003")
	log.Fatal(srv.ListenAndServe())
}
