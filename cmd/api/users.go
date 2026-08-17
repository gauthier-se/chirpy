package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gauthier-se/chirpy/internal/auth"
	"github.com/gauthier-se/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
		return
	}

	err := cfg.DB.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Couldn't delete users")
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't hash password")
		return
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 500, "couldn't create user")
		return
	}

	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	responseWithJSON(w, 201, responseUser)
}
