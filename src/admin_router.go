package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getAdminRouter(apiCfg *apiConfig) chi.Router {
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.serveAdminMetrics)
	return adminRouter
}

func (cfg *apiConfig) serveAdminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	html := `<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %v times!</p>
	</body>
	
	</html>
	`
	fmt.Fprintf(w, html, cfg.numServerHits)
}
