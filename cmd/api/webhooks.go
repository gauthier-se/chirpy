package main

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gauthier-se/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || subtle.ConstantTimeCompare([]byte(apiKey), []byte(cfg.PolkaKey)) != 1 {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.DB.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
