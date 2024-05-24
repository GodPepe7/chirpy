package db

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type RefreshToken struct {
	Token     string    `json:"refresh_token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserId    int       `json:"user_id"`
}

func (db *DB) CreateRefreshToken(userId int) (RefreshToken, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("error generating 32 bytes random data")
	}
	token := hex.EncodeToString(randomBytes)
	expiresAt := time.Now().Add(time.Duration(60 * 24 * time.Hour))
	refreshToken := RefreshToken{Token: token, ExpiresAt: expiresAt, UserId: userId}
	dbStruct.Tokens[userId] = refreshToken
	db.writeDB(dbStruct)
	return refreshToken, nil
}

func (db *DB) GetRefreshTokenByToken(token string) (RefreshToken, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	for key, refreshToken := range dbStruct.Tokens {
		if refreshToken.Token == token {
			if time.Now().Compare(refreshToken.ExpiresAt) != -1 {
				delete(dbStruct.Tokens, key)
				return RefreshToken{}, fmt.Errorf("refresh token is expired")
			}
			return refreshToken, nil
		}
	}
	return RefreshToken{}, fmt.Errorf("refresh token doesn't exist with token: %v", token)
}

func (db *DB) DeleteRefreshToken(token string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	for key, refreshToken := range dbStruct.Tokens {
		if refreshToken.Token == token {
			delete(dbStruct.Tokens, key)
			db.writeDB(dbStruct)
			return nil
		}
	}
	return fmt.Errorf("refresh token found not found")
}
