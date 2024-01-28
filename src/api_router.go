package main

import (
	"github.com/go-chi/chi/v5"
)

func getApiRouter(apiCfg *apiConfig) chi.Router {
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", checkReadiness)
	apiRouter.Get("/metrics", apiCfg.serveMetrics)
	apiRouter.Get("/reset", apiCfg.resetMetrics)

	apiRouter.Post("/chirps", apiCfg.createChirp)
	apiRouter.Get("/chirps", apiCfg.getChirps)
	apiRouter.Get("/chirps/{ID}", apiCfg.getChirpByID)
	apiRouter.Delete("/chirps/{ID}", apiCfg.deleteChirp)

	apiRouter.Post("/users", apiCfg.createUser)
	apiRouter.Get("/users", apiCfg.getUsers)
	apiRouter.Get("/users/{ID}", apiCfg.getUserByID)
	apiRouter.Post("/login", apiCfg.loginUser)
	apiRouter.Put("/users", apiCfg.updateUser)

	apiRouter.Post("/refresh", apiCfg.refreshToken)
	apiRouter.Post("/revoke", apiCfg.revokeToken)

	apiRouter.Post("/polka/webhooks", apiCfg.upgradeUser)

	return apiRouter
}
