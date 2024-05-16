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
	id   int    `json:"id"`
	body string `json:"body"`
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
	content, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, fmt.Errorf("error while reading db file: %v", err)
	}
	dbStruct := DBStructure{}
	err = json.Unmarshal(content, &dbStruct)
	if err != nil {
		return DBStructure{}, fmt.Errorf("error unmarshaling db file content '%v': %v", string(content), err)
	}
	return dbStruct, nil
}

func (db *DB) writeDB(dbStruct DBStructure) error {
	jsonByte, err := json.Marshal(dbStruct.Chirps)
	if err != nil {
		return fmt.Errorf("error while marshaling db structure %v: %v", dbStruct.Chirps, err)
	}
	err = os.WriteFile(db.path, jsonByte, 0666)
	if err != nil {
		return fmt.Errorf("error while writing '%v' to file '%v': %v", string(jsonByte), db.path, err)
	}
	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{id, body}
	dbStruct.Chirps[id] = chirp
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, fmt.Errorf("error writing %v to db file", dbStruct.Chirps)
	}
	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
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
