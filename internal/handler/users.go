package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/godpepe7/chirpy/internal/db"
	"github.com/godpepe7/chirpy/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *ApiConfig) PostUserHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := UserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	if existingUser, err := cfg.DB.GetUserByEmail(params.Email); existingUser.Id != 0 && err == nil {
		utils.RespondWithError(rw, 400, fmt.Sprintf("user already exists with email: %v", params.Email))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong creating the user")
	}
	user, err := cfg.DB.CreateUser(params.Email, string(hashedPassword))
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong creating the user")
		return
	}
	userResponse := UserResponse{Id: user.Id, Email: user.Email, IsChirpyRed: user.IsChirpyRed}
	utils.RespondWithJSON(rw, 201, userResponse)
}

func (cfg *ApiConfig) PutUserHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := UserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	jwt, err := utils.GetTokenFromHeader(req, cfg.JwtSecret)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 401, "Invalid header format or token")
		return
	}
	idString, err := jwt.Claims.GetSubject()
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Error getting id from token")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong updating the user")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong updating the user")
		return
	}
	user, err := cfg.DB.GetUserById(id)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 404, "User doesn't exist")
		return
	}
	updatedUser, err := cfg.DB.UpdateUser(id, db.UserParams{
		Email:       params.Email,
		Password:    string(hashedPassword),
		IsChirpyRed: user.IsChirpyRed,
	})
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong updating the user")
		return
	}
	utils.RespondWithJSON(rw, 200, UserResponse{
		Id:          updatedUser.Id,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}
