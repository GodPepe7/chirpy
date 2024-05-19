package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/godpepe7/chirpy/internal/utils"
)

type userParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) PostUserHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}

	user, err := h.DB.CreateUser(params.Email)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong creating the user")
		return
	}
	utils.RespondWithJSON(rw, 201, user)
}
