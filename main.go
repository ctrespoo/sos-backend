package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	authm "sos/backend/auth"
	categoria "sos/backend/categorias"
	produtos "sos/backend/produtos"

	interno "sos/backend/interno/db"

	global "sos/backend/global"

	"cloud.google.com/go/storage"
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

	config := &firebase.Config{
		StorageBucket: "sos-do-maceneiro.appspot.com",
	}
	opt := option.WithCredentialsFile("./firebase.json")
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatal(err)
	}
	cliente, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	clientStorage, err := app.Storage(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	bucket, err := clientStorage.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
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
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Post("/", home(cliente, bucket))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authm.Login(dbtx, cliente))
		r.Post("/auth/verificar", authm.Verificar(dbtx, cliente))
		r.Post("/auth/registrar", authm.CriarConta(dbtx, cliente, db))

		r.Group(func(r chi.Router) {
			r.Get("/produtos/{id}", produtos.PegarProdutoUnico(dbtx, cliente, bucket))
			r.Get("/produtos", produtos.PegarTodosProdutos(dbtx, cliente))
			r.Post("/produtos", produtos.CriaProduto(dbtx, cliente, bucket))
			r.Get("/categorias", categoria.PegarTodasCategorias(dbtx, cliente))
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

func home(client *auth.Client, bucket *storage.BucketHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Header.Get("Content-Type"))
		imagem, _, err := r.FormFile("imagem")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao pegar imagem"}`))
			return
		}
		defer imagem.Close()
		imgBytes, _ := io.ReadAll(imagem)
		write := bucket.Object("teste/go.png").NewWriter(r.Context())
		defer write.Close()
		b, err := io.Copy(write, bytes.NewReader(imgBytes))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "Erro ao pegar imagem"}`))
			return
		}

		log.Println(b)

		w.Write([]byte("b"))
	}
}
