package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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
	dbStruct := DBStructure{Chirps: map[int]Chirp{}, Users: map[int]User{}}
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

func RemoveDBFile() {
	filePath := "../../database.json"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}
	e := os.Remove(filePath)
	if e != nil {
		log.Fatal(e)
	}
}
