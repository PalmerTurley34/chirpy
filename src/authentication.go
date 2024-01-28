package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func newToken(issuer string, expiresIn time.Duration, subject string) *jwt.Token {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   subject,
		},
	)
}

func getTokenSignature(token *jwt.Token, secretKey string) (string, error) {
	return token.SignedString([]byte(secretKey))
}

func getTokenFromHeader(r *http.Request, secretKey string) (*jwt.Token, error) {
	authTokenString := r.Header.Get("Authorization")
	authTokenString, _ = strings.CutPrefix(authTokenString, "Bearer ")
	return jwt.ParseWithClaims(
		authTokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
}
