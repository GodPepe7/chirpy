package db

import (
	"fmt"

	"github.com/google/uuid"
)

type Chirp struct {
	Id     string `json:"id"`
	Body   string `json:"body"`
	UserId string `json:"author_id"`
}

func (db *DB) CreateChirp(body string, userId string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := uuid.NewString()
	chirp := Chirp{Id: id, Body: body, UserId: userId}
	dbStruct.Chirps[id] = chirp
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, fmt.Errorf("error writing %v to db file", dbStruct.Chirps)
	}
	return chirp, nil
}

func (db *DB) GetChirpById(id string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbStruct.Chirps[id]
	if !ok {
		return Chirp{}, nil
	}
	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := []Chirp{}
	for _, val := range dbStruct.Chirps {
		chirps = append(chirps, val)
	}
	return chirps, nil
}

func (db *DB) DeleteChirp(id string) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	_, ok := dbStruct.Chirps[id]
	if !ok {
		return fmt.Errorf("chirp doesn't exist with id: %v", id)
	}
	delete(dbStruct.Chirps, id)
	return nil
}
