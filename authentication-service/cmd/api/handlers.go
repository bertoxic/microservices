package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(rw http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(rw, r, &requestPayload)
	log.Println("reading json request in authenticate")
	if err != nil {
		app.errorJSON(rw, err, http.StatusBadRequest)
		log.Println(" error occured in reading json request in authenticate")

		return
	}

	//validate user

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	log.Println("getting user by email from database")

	if err != nil {
		log.Println("error occured in getting user by email from database", err.Error())

		app.errorJSON(rw, errors.New(fmt.Sprintf("invalid credentials : %v", err.Error())), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	log.Println("verifying user  password by email from database")

	if err != nil || !valid {
		log.Println("error in validating user by password from database")

		app.errorJSON(rw, errors.New(fmt.Sprintf("invalid credentials : %v", err.Error())), http.StatusBadRequest)
		return
	}
	//log authentication
	//err = app.logRequest( "authenticate", fmt.Sprintf("%s logged in",user.Email))
	if err != nil {

		app.errorJSON(rw, errors.New(fmt.Sprintf("no logging stuff : %v", err.Error())), http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}
	app.writeJSON(rw, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}

//psql -h localhost -p 5432 -U postgres -W
