package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func RespondWithError(rw http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	rw.WriteHeader(code)
	errRes := errorResponse{Error: msg}
	response, _ := json.Marshal(errRes)
	rw.Write(response)
}

func RespondWithJSON(rw http.ResponseWriter, code int, payload interface{}) {
	rw.WriteHeader(code)
	response, _ := json.Marshal(payload)
	rw.Write(response)
}

func ReplaceBadWords(input string) string {
	profane := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}

	split := strings.Split(input, " ")
	for i, word := range split {
		s := strings.ToLower(word)
		censored, ok := profane[s]
		if ok {
			split[i] = censored
		}
	}
	return strings.Join(split, " ")
}

// Creates Jwt for user which expires in 1 hour if 0 is given as expiresInSeconds
func CreateJwt(expiresInSeconds int, userId string, secret string) (string, error) {
	const defaultExpirationInHours = 1
	expiresIn := time.Duration(defaultExpirationInHours * time.Hour)
	if expiresInSeconds != 0 {
		expiresIn = time.Duration(expiresInSeconds * int(time.Second))
	}
	issuedAt := jwt.NewNumericDate(time.Now())
	expiredAt := jwt.NewNumericDate(issuedAt.Add(expiresIn))
	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  issuedAt,
			ExpiresAt: expiredAt,
			Subject:   userId,
		})
	signedJwt, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error while creating jwt: %v", err)
	}
	return signedJwt, nil
}

func ParseJwt(token, secret string) (*jwt.Token, error) {
	jwt, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !jwt.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return jwt, nil
}

func GetTokenFromHeader(req *http.Request, secret string) (*jwt.Token, error) {
	token := req.Header.Get("Authorization")
	if token == "" {
		return nil, fmt.Errorf("no Authorization header found")
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	jwt, err := ParseJwt(token, secret)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse jwt: %v", err)
	}
	return jwt, nil
}
