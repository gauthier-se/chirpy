package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gauthier-se/chirpy/internal/auth"
	"github.com/gauthier-se/chirpy/internal/database"
)

const (
	accessTokenExpiration  = time.Hour
	refreshTokenExpiration = 60 * 24 * time.Hour
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	token, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, accessTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	refreshToken, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpiration),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
		return
	}

	responseWithJSON(w, http.StatusOK, response{
		User:         toUser(dbUser),
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	dbUser, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, accessTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	responseWithJSON(w, http.StatusOK, response{Token: token})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
