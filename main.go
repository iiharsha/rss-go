package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/iiharsha/rss-go/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load()

	port := os.Getenv("PORT") //getting port from .env
	if port == "" {
		log.Fatal("PORT is missing idiot")
	}

	dbUrl := os.Getenv("DB_URL") //getting db_url from .env
	if dbUrl == "" {
		log.Fatal("DB_URL is missing idiot")
	}

	conn, err := sql.Open("postgres", dbUrl) //opening a db connection
	if err != nil {
		log.Fatal("can't connect to database", err)
	}

	//converting conn -> db.Queries
	queries := database.New(conn)

	apiCfg := apiConfig{
		DB: queries,
	}

	router := chi.NewRouter()
	//warning this is like not secure for prod
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.handlerGetUser)

	router.Mount("/v1", v1Router)

	srv := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("server starting on port %v", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("error occured", err)
	}
}
