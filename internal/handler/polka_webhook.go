package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/godpepe7/chirpy/internal/db"
	"github.com/godpepe7/chirpy/internal/utils"
)

type EventData struct {
	UserId int `json:"user_id"`
}

type PolkaWebhookParams struct {
	Event string    `json:"event"`
	Data  EventData `json:"data"`
}

func (cfg ApiConfig) PostPolkaWebhookHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := req.Header.Get("Authorization")
	if token == "" {
		utils.RespondWithError(rw, 401, "No authorization header found")
		return
	}
	apiKey := strings.Replace(token, "ApiKey ", "", 1)
	println(apiKey, cfg.PolkaKey)
	if cfg.PolkaKey != apiKey {
		utils.RespondWithError(rw, 401, "Not authorized for this action")
		return
	}
	decoder := json.NewDecoder(req.Body)
	params := PolkaWebhookParams{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	if params.Event == "user.upgraded" {
		user, err := cfg.DB.GetUserById(params.Data.UserId)
		if err != nil {
			fmt.Println(err)
			utils.RespondWithError(rw, 404, "No such user")
			return
		}
		cfg.DB.UpdateUser(user.Id, db.UserParams{Email: user.Email, Password: user.Password, IsChirpyRed: true})
	}
	rw.WriteHeader(204)
}
