package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/bertoxic/broker/events"
	"github.com/bertoxic/broker/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	log.Println("broker hnaler triggered")
	jax := &jsonResponse{
		Error:   false,
		Data:    "stipend",
		Message: "n",
	}

	err := app.writeJSON(w, http.StatusOK, jax)
	if err != nil {
		log.Println("error in broker handler", err)
		return
	}
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayLoad  `json:"log,omitempty"`
	Mail   MailPayload `jon:"mail,omitempty"`
}

type LogPayLoad struct {
	Name    string `json:"name"`
	Data    string `json:"data"`
	Message string `json:"message,omitempty"`
}

type MailPayload struct {
	FROM    string `json:"from"`
	TO      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) HandleSubmission(rw http.ResponseWriter, r *http.Request) {

	var Requestpayload RequestPayload
	// err := app.readJSON(rw, r, &Requestpayload)
	// if err != nil {
	//     return
	// }
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&Requestpayload)
	if err != nil {
		log.Printf("Error decoding in handlesubmission handler: %s", err.Error())
		app.errorJSON(rw, err)

	}
	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		log.Printf("Error decoding in handlesubmission handler: %s", err.Error())
	}

	switch Requestpayload.Action {
	case "auth":
		app.authenticate(rw, Requestpayload.Auth)
	case "log":
		log.Printf("it is about toooo loggxx oh now: %s", Requestpayload.Log)
		app.LogItemRPC(rw, Requestpayload.Log)
		//app.LogEventViaRabbit(rw, Requestpayload.Log)
		log.Printf("it is loggxx oh now: %s", Requestpayload.Log)
		//app.logItem(rw, Requestpayload.Log)

	case "mail":
		app.sendMail(rw, Requestpayload.Mail)
	default:
		app.errorJSON(rw, errors.New("unknown action"))
		log.Printf("it is about toooo loggxx oh now unkonwn action broker handler: %s", Requestpayload)
	}
}

func (app *Config) sendMail(rw http.ResponseWriter, mail MailPayload) {
	const (
		mailURL = "http://mail-service/send"
	)
	jsonData, err := json.MarshalIndent(mail, "", "\t")
	if err != nil {
		println("error in marshalling data in mailsend", err)
	}

	request, err := http.NewRequest(http.MethodPost, mailURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		app.errorJSON(rw, fmt.Errorf("error in making request to mail /send route: %s", err.Error()))
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(rw, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Println("erooor in broker satuscode rtrunb", err.Error())
		app.errorJSON(rw, fmt.Errorf("error calling mail service having statuscode of : %d", response.StatusCode))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Data = ""
	payload.Message = "Message sent to" + mail.TO
	err = app.writeJSON(rw, http.StatusAccepted, payload)
	if err != nil {
		log.Println("erooor in broker satuscode", err.Error())
		app.errorJSON(rw, err)
		return

	}

}

func (app *Config) logItem(rw http.ResponseWriter, entry LogPayLoad) {
	//createjson to send to auth microservice
	log.Println("data gotten is", entry.Name, entry.Data)
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		println("Error marshalling data zzzzzzzzzzzzzz:", err)
		return
	}
	logServiceURL := "http://logger-service/log"
	println("gotten to request now wooow", jsonData)
	//call the service
	request, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		println("eroooosss in making request in boker handler vvv:")
		app.errorJSON(rw, err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		println("eroooosss in sending through client oh respone vvv:", err.Error())
		app.errorJSON(rw, err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		println("Error reading response body:", body, err)
		// Handle the error appropriately (e.g., return an error to the caller)
		return
	}
	if response != nil {
		defer response.Body.Close()
	}
	if response.StatusCode != http.StatusAccepted {

		app.errorJSON(rw, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged successfully"

	app.writeJSON(rw, http.StatusAccepted, payload)

}
func (app *Config) authenticate(rw http.ResponseWriter, a AuthPayload) {
	//createjson to send to auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the service
	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(rw, err)

	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(rw, err)
		return
	}
	defer response.Body.Close()

	//makes sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(rw, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(rw, fmt.Errorf("somthing up , error calling auth service error: %s", response.Body))
		return

	}

	//create variable to read response.boddy into

	var jsonFromService jsonResponse
	// decode the json from the auth
	err = json.NewDecoder(request.Body).Decode(&jsonFromService)
	if jsonFromService.Error {
		app.errorJSON(rw, errors.Join(err, errors.New("just check the previous error biko")), http.StatusUnauthorized)
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Message

	app.writeJSON(rw, http.StatusAccepted, payload)
}

func (app *Config) LogEventViaRabbit(rw http.ResponseWriter, l LogPayLoad) {
	log.Printf("at event logger  currently %s,  %s", l.Name, l.Data)
	err := app.PushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(rw, fmt.Errorf("error in pushqueue: %s", err.Error()))
		return
	}
	log.Printf("done with event logger  currently %s", l)
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via rabbitmq"
	payload.Data = "logged... via rabbitmq data"
	log.Printf("sending this payload now %s.., %s", payload.Data, payload.Error)

	app.writeJSON(rw, http.StatusAccepted, payload)

}

func (app *Config) PushToQueue(name, msg string) error {
	log.Printf("pushing to queue currently")
	emitter, err := events.NewEmittter(app.Rabbit)
	if err != nil {
		log.Printf("error in broker main function pushTOQueue* : %s", err.Error())
		return err
	}

	payLoad := LogPayLoad{
		Name: name,
		Data: msg,
	}

	j, err := json.MarshalIndent(payLoad, "", "\t")
	if err != nil {
		log.Printf("error in broker main function pushTOQueue* 2: %s", err.Error())
		return err
	}
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		log.Printf("error in broker main function pushTOQueue* 3: %s", err.Error())
		return err
	}

	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) LogItemRPC(rw http.ResponseWriter, l LogPayLoad) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		log.Printf("error in logitemrpc currently tring to connect to client:%s", err.Error())
		app.errorJSON(rw, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}
	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		log.Printf("error in logitemrpc currently %s", err.Error())
		app.errorJSON(rw, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Data:    "sennt duccessfully",
		Message: fmt.Sprintf("data sent was %s", l),
	}

	app.writeJSON(rw, http.StatusAccepted, payload)
}

func (app *Config) LogviaGRPC(rw http.ResponseWriter, r *http.Request) {
	log.Println("inside log via grpc in handler")

	var requestPayload RequestPayload
	err := app.readJSON(rw, r, &requestPayload)
	if err != nil {
		log.Printf("error occured while reading json in handlers is %s", err.Error())
		app.errorJSON(rw, err)
		return
	}

	log.Println("inside log via grpc in making conn")
	var conn *grpc.ClientConn
	conn, err = grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Printf("error cant make connection %s", err.Error())
		app.errorJSON(rw, err)
		return
	}
	defer conn.Close()

	log.Println("connection made, about making client")
	c := logs.NewLogServiceClient(conn)
	log.Println("client made, about calling write log in broker handler")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("about to call write log in broker")
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		log.Printf("error occured which is %s", err.Error())
		app.errorJSON(rw, err)
		return
	} else {
		log.Println("WriteLog call successful")
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged in message"
	payload.Data = "logged in data"

	app.writeJSON(rw, http.StatusAccepted, payload)

}
