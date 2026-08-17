package main

import "net/http"

func (cfg *apiConfig) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", healthHandler)

	mux.HandleFunc("POST /api/users", cfg.createUserHandler)

	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)

	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)

	return mux
}
