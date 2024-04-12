package main

import (
	"net/http"

	"github.com/bertoxic/log-service/data"
)



type JSONPayload struct {
    Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(rw http.ResponseWriter, r *http.Request) {
    // ead json into var
    var requestPayLoad JSONPayload
    err := app.readJSON( rw, r, &requestPayLoad)
    if err != nil {
        println("tooomany error logs hhhhhhhhh:")
		app.errorJSON(rw, err)
	}

    event := data.LogEntry {
        Name: requestPayLoad.Name,
        Data: requestPayLoad.Data,
    }
    println("rexsffgfgfgadponzx:",event.Data,event.Name)
    err = app.Models.LogEntry.Insert(event)
    if err != nil {
        println("errors in writelog in logger:")

        app.errorJSON(rw, err)
        return
    }
    resp := jsonResponse{
        Error: false,
        Message: "logged fully no issues",
        //Data: "finally logged us in right",
    }
    app.writeJSON(rw, http.StatusAccepted, resp)
} 