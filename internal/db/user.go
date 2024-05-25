package db

import (
	"fmt"
)

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserParams struct {
	Email       string
	Password    string
	IsChirpyRed bool
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStruct.Users) + 1
	user := User{Id: id, Email: email, Password: password, IsChirpyRed: false}
	dbStruct.Users[id] = user
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, fmt.Errorf("error writing %v to db file", dbStruct.Users)
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, nil
}

func (db *DB) UpdateUser(id int, updated UserParams) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStruct.Users[id]
	if !ok {
		return User{}, fmt.Errorf("no such user: %v", id)
	}
	dbStruct.Users[id] = User{
		Id:          user.Id,
		Email:       updated.Email,
		Password:    updated.Password,
		IsChirpyRed: updated.IsChirpyRed,
	}
	db.writeDB(dbStruct)
	return dbStruct.Users[id], nil
}

func (db *DB) GetUserById(id int) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStruct.Users[id]
	if !ok {
		return User{}, fmt.Errorf("no such user: %v", id)
	}
	return user, nil
}
