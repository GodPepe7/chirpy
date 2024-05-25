package db

import (
	"fmt"
)

type Chirp struct {
	Id     int    `json:"id"`
	Body   string `json:"body"`
	UserId int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, userId int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{Id: id, Body: body, UserId: userId}
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
		return Chirp{}, fmt.Errorf("no chirp with id: %v", id)
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

func (db *DB) DeleteChirp(id int) error {
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
	err = db.writeDB(dbStruct)
	if err != nil {
		return fmt.Errorf("error writing %v to db file", dbStruct.Chirps)
	}
	return nil
}
