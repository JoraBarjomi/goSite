package main

import (
	"net/http"
	handlers "site/server"
	utils "site/server/util"
)

func main() {
	http.HandleFunc("/", handlers.MainHandler)
	http.HandleFunc("/notes", handlers.MainHandler)
	http.HandleFunc("/notes/new", handlers.NewHandler)
	http.HandleFunc("/notes/create", handlers.CreateHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	utils.Check(err)
}
