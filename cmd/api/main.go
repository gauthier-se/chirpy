package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gauthier-se/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB:        dbQueries,
		Platform:  platform,
		JWTSecret: jwtSecret,
	}

	mux := apiCfg.setupRoutes()

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Serveur démarré sur le port %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
