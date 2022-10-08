package models

import (
	"context"
	"go-rest/db"
	"log"
)

type User struct {
	Id, Name, Email, Password, ProfilePicture, Biography string
}

type create interface {
	CreateUser(id string, name string, email string, password string) error
}

type read interface {
	CheckExistingEmail(email string) map[string]interface{}
	FindById(id string) User
}

func (user User) CreateUser(id string, name string, email string, password string) error {
	queryString := `
		INSERT INTO users(id, name, email, password) VALUES($1,$2,$3,$4)
	`

	Conn := db.Conn

	_,insertErr := Conn.Exec(context.Background(),queryString,id,name,email,password);
	if insertErr != nil {
		return insertErr
	}
	return nil
}

func (user User) CheckExistingEmail(email string) map[string]interface{} {
	queryString := `
	SELECT name,email,password,profile_picture,biography FROM users WHERE email = $1
	`

	Conn := db.Conn
	data,dataErr := Conn.Query(context.Background(),queryString, email)

	if dataErr != nil {
		log.Fatal(dataErr)
	}

	var userResult User
	for data.Next() {
		data.Scan(&userResult.Name,&userResult.Email,&userResult.Password,&userResult.ProfilePicture,&userResult.Biography)
	}

	resp := make(map[string]interface{})

	if userResult.Email != "" {
		resp["isExist"] = true
		resp["user"] = userResult
		return resp
	} else {
		resp["isExist"] = false
		resp["user"] = userResult
		return resp
	}
}

func (user User) FindById(id string) User {
	queryString := `
		SELECT name,email,password,profile_picture,biography FROM users WHERE id = $1
	`

	Conn := db.Conn
	data,dataErr := Conn.Query(context.Background(),queryString, id)

	if dataErr != nil {
		log.Fatal(dataErr)
	}

	var userData = User{}

	for data.Next() {
		scanErr := data.Scan(&userData.Name, &userData.Email, &userData.Password, &userData.ProfilePicture, &userData.Biography)

		if scanErr != nil {
			log.Fatal(scanErr)
		}
	}

	return userData

}