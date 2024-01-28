package main

import (
	"fmt"
	"net/http"
	"time"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(r, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "invalid refresh token")
		return
	}
	issuer, _ := token.Claims.GetIssuer()
	if issuer != "chirpy-refresh" {
		respondWithError(w, 401, "authorization token must be refresh token, not access token")
		return
	}
	dbToken, err := cfg.database.GetToken(token.Raw)
	if err != nil {
		respondWithError(w, 401, "token is not recognized")
		return
	}
	if dbToken.IsRevoked {
		respondWithError(w, 401, "token has been revoked")
		return
	}
	accessToken := newToken("chirpy-access", time.Hour, fmt.Sprint(dbToken.UserID))
	signedAccessToken, err := getTokenSignature(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating access token: %v", err.Error()))
		return
	}
	type tokenResponse struct {
		Token string `json:"token"`
	}
	respondWithJson(w, 200, tokenResponse{Token: signedAccessToken})
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(r, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "invalid refresh token")
		return
	}
	issuer, _ := token.Claims.GetIssuer()
	if issuer != "chirpy-refresh" {
		respondWithError(w, 401, "authorization token must be refresh token, not access token")
		return
	}
	_, err = cfg.database.RevokeToken(token.Raw)
	if err != nil {
		respondWithError(w, 401, "token is not recognized")
		return
	}
	w.WriteHeader(200)
}
