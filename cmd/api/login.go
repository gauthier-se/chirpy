package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gauthier-se/chirpy/internal/auth"
)

const defaultExpiration = time.Hour

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expiresIn := defaultExpiration
	if params.ExpiresInSeconds > 0 {
		expiresIn = min(time.Duration(params.ExpiresInSeconds)*time.Second, defaultExpiration)
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}

	responseWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token: token,
	})
}
