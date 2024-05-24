package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/godpepe7/chirpy/internal/utils"
)

type chirpParams struct {
	Body string `json:"body"`
}

func (cfg *ApiConfig) GetChirpHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong with getting chirps")
		return
	}
	utils.RespondWithJSON(rw, 200, chirps)
}

func (cfg *ApiConfig) GetChirpByIdHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	chirpId := req.PathValue("id")
	chirp, err := cfg.DB.GetChirpById(chirpId)
	if chirp.Id == "" && err == nil {
		utils.RespondWithError(rw, 404, fmt.Sprintf("Chirp with ID %v doesn't exist", chirpId))
		return
	}
	utils.RespondWithJSON(rw, 200, chirp)
}

func (cfg *ApiConfig) PostChirpHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := chirpParams{}
	err := decoder.Decode(&params)

	if err != nil {
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		utils.RespondWithError(rw, 400, "Chirp is too long")
		return
	}
	jwt, err := utils.GetTokenFromHeader(req, cfg.JwtSecret)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 401, "Invalid token")
		return
	}
	userId, err := jwt.Claims.GetSubject()
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Error getting id from token")
		return
	}
	cleanedString := utils.ReplaceBadWords(params.Body)
	chirp, err := cfg.DB.CreateChirp(cleanedString, userId)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	utils.RespondWithJSON(rw, 201, chirp)
}
