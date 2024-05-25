package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/godpepe7/chirpy/internal/db"
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
	authorFilter := req.URL.Query().Get("author_id")
	sortFilter := req.URL.Query().Get("sort")
	if sortFilter == "" || sortFilter == "asc" {
		sort.SliceStable(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})
	} else if sortFilter == "desc" {
		sort.SliceStable(chirps, func(i, j int) bool {
			return chirps[i].Id > chirps[j].Id
		})
	}
	if authorFilter == "" {
		utils.RespondWithJSON(rw, 200, chirps)
		return
	}
	userId, err := strconv.Atoi(sortFilter)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 401, "Pls enter a valid id number value")
	}
	filtered := []db.Chirp{}
	for _, chirp := range chirps {
		if chirp.UserId == userId {
			filtered = append(filtered, chirp)
		}
	}
	utils.RespondWithJSON(rw, 200, filtered)
}

func (cfg *ApiConfig) GetChirpByIdHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	idValue := req.PathValue("id")
	chirpID, err := strconv.Atoi(idValue)
	if err != nil {
		utils.RespondWithError(rw, 400, "Something went wrong parsing the id, has to be a number")
		return
	}
	chirp, err := cfg.DB.GetChirpById(chirpID)
	if chirp.Id == 0 && err == nil {
		utils.RespondWithError(rw, 404, fmt.Sprintf("Chirp with ID %v doesn't exist", chirpID))
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
		utils.RespondWithError(rw, 401, "Invalid header format or token")
		return
	}
	userIdString, err := jwt.Claims.GetSubject()
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Error getting id from token")
		return
	}
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong updating the user")
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

func (cfg *ApiConfig) DeleteChirpHandler(rw http.ResponseWriter, req *http.Request) {
	chirpIdString := req.PathValue("id")
	jwt, err := utils.GetTokenFromHeader(req, cfg.JwtSecret)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 401, "Invalid token header")
		return
	}
	chirpId, err := strconv.Atoi(chirpIdString)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong parsing the id, has to be a number")
		return
	}
	userIdString, err := jwt.Claims.GetSubject()
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Error getting id from token")
		return
	}
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 500, "Something went wrong parsing the id, has to be a number")
		return
	}
	println(userId)
	chirp, err := cfg.DB.GetChirpById(chirpId)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 404, fmt.Sprintf("Chirp with ID %v couldn't be found", chirpId))
		return
	}
	fmt.Println(chirp)
	if chirp.UserId != userId {
		utils.RespondWithError(rw, 403, "Unauthorized")
		return
	}
	err = cfg.DB.DeleteChirp(chirpId)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(rw, 401, fmt.Sprintf("No such chirp with id: %v", chirpId))
		return
	}
}
