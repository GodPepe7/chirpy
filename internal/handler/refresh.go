package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/godpepe7/chirpy/internal/utils"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) PostRefreshHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	headerToken := req.Header.Get("Authorization")
	if headerToken == "" {
		utils.RespondWithError(rw, 400, "Need authorization header")
		return
	}
	headerToken = strings.Replace(headerToken, "Bearer ", "", 1)
	refreshToken, err := cfg.DB.GetRefreshTokenByToken(headerToken)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 404, fmt.Sprintf("Couldn't find Refresh Token '%v'", refreshToken))
		return
	}
	accessToken, err := utils.CreateJwt(0, refreshToken.UserId, cfg.JwtSecret)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong creating new access token")
		return
	}
	utils.RespondWithJSON(rw, 200, RefreshResponse{Token: accessToken})
}
