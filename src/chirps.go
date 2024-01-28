package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PalmerTurley34/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	token, err := getTokenFromHeader(r, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Error with auth token")
		return
	}
	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, 401, "auth token has no user ID")
	}
	userID, _ := strconv.Atoi(userIDStr)

	decoder := json.NewDecoder(r.Body)
	rBody := requestBody{}
	err = decoder.Decode(&rBody)
	if err != nil {
		respondWithError(w, 500, "Could not decode request body")
		return
	}
	const maxChirpLen = 140
	if len(rBody.Body) > maxChirpLen {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	newChirp, err := cfg.database.CreateChirp(userID, replaceBadWords(rBody.Body))
	if err != nil {
		respondWithError(w, 500, "Error creating new chirp")
		return
	}
	respondWithJson(w, 201, newChirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStirng := r.URL.Query().Get("author_id")
	sortBy := r.URL.Query().Get("sort")
	var authorID int
	var err error
	if authorIDStirng == "" {
		authorID = -1
	} else {
		authorID, err = strconv.Atoi(authorIDStirng)
		if err != nil {
			respondWithError(w, 400, "author ID must be an interger")
			return
		}
	}
	allChirps, err := cfg.database.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Error fetching Chirps")
		return
	}
	chirpsSlice := make([]database.Chirp, 0, len(allChirps))
	for _, chirp := range allChirps {
		if authorID > 0 && chirp.UserID != authorID {
			continue
		}
		chirpsSlice = append(chirpsSlice, chirp)
	}
	sortFunc := func(i, j int) bool { return chirpsSlice[i].ID < chirpsSlice[j].ID }
	if sortBy == "desc" {
		sortFunc = func(i, j int) bool { return chirpsSlice[i].ID > chirpsSlice[j].ID }
	}
	sort.Slice(chirpsSlice, sortFunc)
	respondWithJson(w, 200, chirpsSlice)
}

func (cfg *apiConfig) getChirpByID(w http.ResponseWriter, r *http.Request) {
	allChirps, err := cfg.database.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Error fetching Chirps")
		return
	}
	idStr := chi.URLParam(r, "ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, 400, "ID must be an integer")
		return
	}
	chirp, ok := allChirps[id]
	if !ok {
		respondWithError(w, 404, "ID Not Found")
		return
	}
	respondWithJson(w, 200, chirp)
}

func replaceBadWords(msg string) string {
	const (
		badWord1          = "kerfuffle"
		badWord2          = "sharbert"
		badWord3          = "fornax"
		badWordReplacemnt = "****"
	)
	words := strings.Split(msg, " ")
	filterdWords := make([]string, len(words))
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if lowerWord == badWord1 || lowerWord == badWord2 || lowerWord == badWord3 {
			filterdWords[i] = badWordReplacemnt
		} else {
			filterdWords[i] = word
		}
	}
	return strings.Join(filterdWords, " ")
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "ID")
	chirpID, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, 400, "ID must be an integer")
		return
	}
	token, err := getTokenFromHeader(r, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Error with auth token")
		return
	}
	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, 401, "auth token has no user ID")
	}
	userID, _ := strconv.Atoi(userIDStr)
	chirps, _ := cfg.database.GetChirps()
	chirpToDelete, ok := chirps[chirpID]
	if !ok {
		respondWithError(w, 400, "chirp does not exist")
		return
	}
	if userID != chirpToDelete.ID {
		respondWithError(w, 403, "not authorized to delte chirp")
		return
	}
	cfg.database.DeleteChirp(chirpID)
	w.WriteHeader(200)
}
