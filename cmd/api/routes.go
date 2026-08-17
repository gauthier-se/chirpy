package main

import "net/http"

func (cfg *apiConfig) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", healthHandler)

	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("PUT /api/users", cfg.updateUserHandler)

	mux.HandleFunc("POST /api/login", cfg.loginHandler)
	mux.HandleFunc("POST /api/refresh", cfg.refreshHandler)
	mux.HandleFunc("POST /api/revoke", cfg.revokeHandler)

	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirpHandler)

	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)

	return mux
}
