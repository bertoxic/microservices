package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) SendMail(rw http.ResponseWriter, r *http.Request) {
    type mailMessage struct {
        From string `json:"from"`
        To string `json:"to"`
        Subject string `json:"subject"`
        Message string `json:"message"`
    }

    var requestPayLoad mailMessage
    err := app.readJSON(rw, r, &requestPayLoad)
    if err != nil {
        app.errorJSON(rw, err)
        log.Println("error occured in mailservice 2:",err)

        return
    }
    msg := Message {
        From: requestPayLoad.From,
        To: requestPayLoad.To,
        Subject: requestPayLoad.Subject,
        Data: requestPayLoad.Message,
    }

    err = app.Mailer.SendSMPTMessage(msg)
    if err != nil {
        app.errorJSON(rw, err)
        log.Println("error occured in mailservice 1:",err)
        return
    }
    log.Println("response sent is :",requestPayLoad)

    payload := jsonResponse {
        Error: false,
        Data: "email message sent successfully ",
        Message: fmt.Sprintf("obtained message ohhhh"),
    }

    app.writeJSON(rw, http.StatusAccepted,payload)
}

