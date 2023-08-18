package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	authm "sos/backend/auth"
	produtos "sos/backend/produtos"

	interno "sos/backend/interno/db"

	global "sos/backend/global"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/option"
)

func init() {
	global.TokenAuth = jwtauth.New("HS256", []byte("secret"), nil) // replace with secret key
}

func main() {
	ctx := context.Background()

	dbs := os.Getenv("DATABASE_URL")
	if dbs == "" {
		log.Fatal("DATABASE_URL não definida")
	}

	db, err := pgxpool.New(ctx, dbs)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dbtx := interno.New(db)

	opt := option.WithCredentialsFile("./firebase.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal(err)
	}
	cliente, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", home(cliente))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authm.Login(dbtx, cliente))
		r.Post("/auth/verificar", authm.Verificar(dbtx, cliente))
		r.Post("/auth/registrar", authm.CriarConta(dbtx, cliente, db))

		r.Group(func(r chi.Router) {
			r.Get("/produtos", produtos.PegarTodosProdutos(dbtx, cliente))
			r.Post("/produtos", produtos.CriaProduto(dbtx, cliente))
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

func home(client *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := "f866ba49-e72a-4ec0-90c5-e3d014722be4"

		claims := map[string]interface{}{
			"premiumAccount": true,
		}

		err := client.SetCustomUserClaims(r.Context(), uid, claims)
		if err != nil {
			log.Fatalf("error setting custom claims %v\n", err)
		}
		u, err := client.GetUser(r.Context(), uid)
		if err != nil {
			log.Fatalf("error getting user %s: %v\n", uid, err)
		}

		token, err := client.CustomTokenWithClaims(r.Context(), uid, claims)
		if err != nil {
			log.Fatalf("error minting custom token: %v\n", err)
		}

		log.Printf("Successfully fetched user data: %v\n", u.CustomClaims)
		w.Write([]byte(token))
	}
}
