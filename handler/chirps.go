package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/godpepe7/chirpy/internal/utils"
)

type chirpParams struct {
	Body string `json:"body"`
}

func (h *Handler) GetChirpHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	chirps, err := h.DB.GetChirps()
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong with getting chirps")
		return
	}
	utils.RespondWithJSON(rw, 200, chirps)
}

func (h *Handler) GetChirpByIdHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	idValue := req.PathValue("id")
	chirpID, err := strconv.Atoi(idValue)
	if err != nil {
		utils.RespondWithError(rw, 400, "Something went wrong parsing the id, has to be a number")
		return
	}
	chirp, err := h.DB.GetChirpById(chirpID)
	if err != nil {
		utils.RespondWithError(rw, 404, fmt.Sprintf("Chirp with ID %v doesn't exist", chirpID))
		return
	}
	utils.RespondWithJSON(rw, 200, chirp)
}

func (h *Handler) PostChirpHandler(rw http.ResponseWriter, req *http.Request) {
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

	cleanedString := utils.ReplaceBadWords(params.Body)
	chirp, err := h.DB.CreateChirp(cleanedString)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong")
		return
	}
	utils.RespondWithJSON(rw, 201, chirp)
}
