package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	Token     string    `json:"refresh_token"`
	UserID    int       `json:"user_id"`
	IsRevoked bool      `json:"is_revoked"`
	RevokedAt time.Time `json:"revoked_at,omitempty"`
}

func (db *DB) CreateToken(token string, userID int) (RefreshToken, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	if _, ok := dbStruct.RefreshTokens[token]; ok {
		return RefreshToken{}, errors.New("token already exists")
	}
	newToken := RefreshToken{Token: token, UserID: userID}
	dbStruct.RefreshTokens[token] = newToken
	err = db.writeDB(dbStruct)
	if err != nil {
		return RefreshToken{}, err
	}
	return newToken, nil
}

func (db *DB) GetToken(token string) (RefreshToken, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	refreshToken, ok := dbStruct.RefreshTokens[token]
	if !ok {
		return RefreshToken{}, errors.New("token does not exist")
	}
	return refreshToken, nil
}

func (db *DB) RevokeToken(token string) (RefreshToken, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	refreshToken, ok := dbStruct.RefreshTokens[token]
	if !ok {
		return RefreshToken{}, errors.New("token does not exist")
	}
	refreshToken.RevokedAt = time.Now().UTC()
	refreshToken.IsRevoked = true
	dbStruct.RefreshTokens[token] = refreshToken
	err = db.writeDB(dbStruct)
	if err != nil {
		return RefreshToken{}, err
	}
	return refreshToken, nil
}
