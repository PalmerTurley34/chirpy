package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	type webhookBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	polkaAuthHeader := r.Header.Get("Authorization")
	polkaAuthHeader, _ = strings.CutPrefix(polkaAuthHeader, "ApiKey ")
	if polkaAuthHeader != cfg.polkaAPIKey {
		respondWithError(w, 401, "invalid API key")
		return
	}
	rBody := webhookBody{}
	err := parseReqBody(r, &rBody)
	if err != nil {
		respondWithError(w, 500, "could not parse body")
	}
	if rBody.Event != "user.upgraded" {
		w.WriteHeader(200)
		return
	}
	users, _ := cfg.database.GetUsers()
	user, ok := users[rBody.Data.UserID]
	if !ok {
		respondWithError(w, 404, "user does not exist")
	}
	cfg.database.UpgradeUserToPremium(user.ID)
	w.WriteHeader(200)
}
