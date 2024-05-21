package db

import (
	"fmt"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserParams struct {
	Email    string
	Password string
}

func (db *DB) CreateUser(email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStruct.Users) + 1
	user := User{Id: id, Email: email, Password: password}
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
	user := dbStruct.Users[id]
	dbStruct.Users[id] = User{Id: user.Id, Email: updated.Email, Password: updated.Password}
	db.writeDB(dbStruct)
	return dbStruct.Users[id], nil
}
