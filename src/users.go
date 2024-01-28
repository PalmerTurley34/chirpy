package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/PalmerTurley34/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type userBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	rBody := userBody{}
	err := decoder.Decode(&rBody)
	if err != nil {
		respondWithError(w, 500, "Could not decode body")
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(rBody.Password), 0)
	if err != nil {
		respondWithError(w, 500, "password cannot be hashed")
		return
	}
	rBody.Password = string(hashedPass)
	newUser, err := cfg.database.CreateUser(rBody.Email, rBody.Password)
	if err != nil {
		respondWithError(w, 500, "Error creating new user")
		return
	}
	respondWithJson(w, 201, removePasswordInResponse(newUser))
}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := cfg.database.GetUsers()
	if err != nil {
		respondWithError(w, 500, "Error fetching Users")
		return
	}
	usersSlice := make([]userResponse, 0, len(allUsers))
	for _, user := range allUsers {
		usersSlice = append(usersSlice, removePasswordInResponse(user))
	}
	sort.Slice(usersSlice, func(i, j int) bool {
		return usersSlice[i].ID < usersSlice[j].ID
	})
	respondWithJson(w, 200, usersSlice)
}

func (cfg *apiConfig) getUserByID(w http.ResponseWriter, r *http.Request) {
	allUsers, err := cfg.database.GetUsers()
	if err != nil {
		respondWithError(w, 500, "Error fetching Users")
		return
	}
	idStr := chi.URLParam(r, "ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, 400, "ID must be an integer")
		return
	}
	user, ok := allUsers[id]
	if !ok {
		respondWithError(w, 404, "ID Not Found")
		return
	}
	respondWithJson(w, 200, removePasswordInResponse(user))
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	rBody := userBody{}
	err := parseReqBody(r, &rBody)
	if err != nil {
		respondWithError(w, 500, "Could not decode body")
		return
	}
	user, err := cfg.database.GetUserByEmail(rBody.Email)
	if err != nil {
		respondWithError(w, 400, "User does not exist")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rBody.Password))
	if err != nil {
		respondWithError(w, 401, "invalid password")
		return
	}
	accessToken := newToken("chirpy-access", time.Hour, fmt.Sprint(user.ID))
	refreshToken := newToken("chirpy-refresh", ((24 * 60) * time.Hour), fmt.Sprint(user.ID))

	returnUser := removePasswordInResponse(user)
	signedAccessToken, err := getTokenSignature(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating access token: %v", err.Error()))
		return
	}
	signedRefreshToken, err := getTokenSignature(refreshToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating refresh token: %v", err.Error()))
		return
	}
	cfg.database.CreateToken(signedRefreshToken, user.ID)
	returnUser.Token = signedAccessToken
	returnUser.RefreshToken = signedRefreshToken
	respondWithJson(w, 200, returnUser)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	rBody := userBody{}
	err := parseReqBody(r, &rBody)
	if err != nil {
		respondWithError(w, 500, "Could not decode body")
		return
	}
	token, err := getTokenFromHeader(r, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("user is not authorized: %v", err.Error()))
		return
	}
	issuer, _ := token.Claims.GetIssuer()
	if issuer != "chirpy-access" {
		respondWithError(w, 401, "authorization token must be access token, not refresh token")
		return
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, 500, "Error validating user")
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(rBody.Password), 0)
	if err != nil {
		respondWithError(w, 500, "password cannot be hashed")
		return
	}
	rBody.Password = string(hashedPass)
	userID, _ := strconv.Atoi(userIDString)
	updatedUser, err := cfg.database.UpdateUser(userID, rBody.Email, rBody.Password)
	if err != nil {
		respondWithError(w, 500, "error updating user")
		return
	}
	respondWithJson(w, 200, removePasswordInResponse(updatedUser))
}

type userResponse struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func removePasswordInResponse(user database.User) userResponse {
	return userResponse{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed}
}
