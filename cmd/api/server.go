package main

import (
	"sync/atomic"

	"github.com/gauthier-se/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	JWTSecret      string
}
