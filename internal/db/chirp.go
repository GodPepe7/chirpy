package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{path: path + ".json", mux: &sync.RWMutex{}}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) ensureDB() error {
	err := os.WriteFile(db.path, []byte(""), 0666)
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	dbStruct := DBStructure{Chirps: map[int]Chirp{}}
	content, err := os.ReadFile(db.path)
	if err != nil {
		return dbStruct, fmt.Errorf("error while reading db file: %v", err)
	}
	if len(content) == 0 {
		return dbStruct, nil
	}
	err = json.Unmarshal(content, &dbStruct)
	if err != nil {
		return dbStruct, fmt.Errorf("error unmarshaling db file content '%v': %v", string(content), err)
	}
	return dbStruct, nil
}

func (db *DB) writeDB(dbStruct DBStructure) error {
	jsonByte, err := json.Marshal(dbStruct)
	if err != nil {
		return fmt.Errorf("error while marshaling db structure %v: %v", dbStruct, err)
	}
	err = os.WriteFile(db.path, jsonByte, 0666)
	if err != nil {
		return fmt.Errorf("error while writing '%v' to file '%v': %v", string(jsonByte), db.path, err)
	}
	return nil
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
