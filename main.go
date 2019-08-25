package main

import (
	"fmt"
	"github.com/dcrvn/dcrvnwww/conf"
	"github.com/dcrvn/dcrvnwww/handlers"
	"github.com/dcrvn/dcrvnwww/models"
	"github.com/go-chi/chi"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting application.")
	conf.Init()
	fmt.Println("Opening database.")
	err := models.Init()
	if err != nil {
		fmt.Println("Database starting was fail.")
		return
	}
	route := chi.NewRouter()
	handlers.Init(route)
	server := generateServer(route)
	server.Addr = conf.PortToServe()
	fmt.Printf("Starting HTTP server on port: %v\n", conf.Port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("httpSrv.ListenAndServe() failed with %s \n", err)
	}
}

func generateServer(route *chi.Mux) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  360 * time.Second,
		Handler:      route,
	}
}