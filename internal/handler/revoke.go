package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/godpepe7/chirpy/internal/utils"
)

func (cfg *ApiConfig) PostRevokeHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	headerToken := req.Header.Get("Authorization")
	if headerToken == "" {
		utils.RespondWithError(rw, 400, "Need authorization header")
		return
	}
	headerToken = strings.Replace(headerToken, "Bearer ", "", 1)
	err := cfg.DB.DeleteRefreshToken(headerToken)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 404, "Response token doesn't exist")
		return
	}
	rw.WriteHeader(204)
}
