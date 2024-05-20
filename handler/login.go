package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/godpepe7/chirpy/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) PostLoginHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := LoginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	user, err := h.db.GetUserByEmail(params.Email)
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
	userResponse := UserResponse{Id: user.Id, Email: user.Email}
	utils.RespondWithJSON(rw, 200, userResponse)
}
