package db

import "fmt"

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStruct.Users) + 1
	user := User{Id: id, Email: email}
	dbStruct.Users[id] = user
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, fmt.Errorf("error writing %v to db file", dbStruct.Users)
	}
	return user, nil
}
