package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/godpepe7/chirpy/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginParams struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type LoginResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *ApiConfig) PostLoginHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := LoginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	if user.Id == 0 {
		fmt.Println(err)
		utils.RespondWithError(rw, 404, fmt.Sprintf("No such user exists with email: %v", params.Email))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		utils.RespondWithError(rw, 401, "Incorrect password")
		return
	}

	const defaultExpirationInHours = 24
	expiresIn := time.Duration(defaultExpirationInHours * time.Hour)
	if params.ExpiresInSeconds != 0 {
		expiresIn = time.Duration(params.ExpiresInSeconds * int(time.Second))
	}
	token, err := utils.CreateJwt(expiresIn, user.Id, cfg.JwtSecret)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong creating the jwt token")
	}

	userResponse := LoginResponse{Id: user.Id, Email: user.Email, Token: token}
	utils.RespondWithJSON(rw, 200, userResponse)
}
