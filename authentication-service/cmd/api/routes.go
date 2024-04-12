package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {

    mux:= chi.NewRouter()

    mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization","Content-Type","X-CSRF-Token","OPTIONS"},
        AllowCredentials: true,
        MaxAge: 300,
    }))

    mux.Use(middleware.Heartbeat("/ping"))

    mux.Post("/authenticate", app.Authenticate)
    return mux
}
//icacls C:\Users\HP\Desktop\Projects\microServices\project\db-data\postgres