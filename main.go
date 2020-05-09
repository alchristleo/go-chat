package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Initialize app...")

	// creating new connection hub
	hub := createConnectionHub()
	go hub.run()

	// serve client build
	fs := http.FileServer(http.Dir("./client/build"))
	http.Handle("/", fs)

	// serve websocket route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(hub, w, r)
	})

	// start server
	var port string

	// let heroku set port in production
	if os.Getenv("GO_ENV") == "PRODUCTION" {
		port = ":" + os.Getenv("PORT")
	} else {
		port = ":8081"
	}

	fmt.Println("Listening to port: ", port)

	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Println("ListenAndServeError: ", err)
	}
}
