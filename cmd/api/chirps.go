package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID  string    `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body      string    `json:"body"`
		UserID  string    `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "something went wrong")
		return
	}

	dbChirp, err := cfg.DB.CreateChirp(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, "couldn't create chirp")
		return
	}

	responseChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp,
		UserID: dbChirp
	}

	responseWithJSON(w, 201, responseChirp)
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "chirp is too long")
		return
	}

	cleanedBody := replaceProfaneWords(params.Body)
	responseWithJSON(w, 200, map[string]string{"cleaned_body": cleanedBody})
}
