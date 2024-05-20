package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/godpepe7/chirpy/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (h *Handler) PostUserHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := UserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	if existingUser, err := h.db.GetUserByEmail(params.Email); existingUser.Id != 0 && err == nil {
		utils.RespondWithError(rw, 400, fmt.Sprintf("user already exists with email: %v", params.Email))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong creating the user")
	}
	user, err := h.db.CreateUser(params.Email, string(hashedPassword))
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong creating the user")
		return
	}
	userResponse := UserResponse{Id: user.Id, Email: user.Email}
	utils.RespondWithJSON(rw, 201, userResponse)
}
