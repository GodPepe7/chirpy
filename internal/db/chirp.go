package db

import (
	"fmt"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{Id: id, Body: body}
	dbStruct.Chirps[id] = chirp
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, fmt.Errorf("error writing %v to db file", dbStruct.Chirps)
	}
	return chirp, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
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
