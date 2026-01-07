package main

import (
	"net/http"
	handlers "site/server"
	utils "site/server/util"
)

func main() {
	http.HandleFunc("/", handlers.MainHandler)
	http.HandleFunc("/notes", handlers.NotesHandler)
	http.HandleFunc("/notes/create", handlers.CreateNoteHandler)
	http.HandleFunc("/todo", handlers.TodoHandler)
	http.HandleFunc("/todo/create", handlers.CreateTodoHandler)
	http.HandleFunc("/todo/update", handlers.UpdateTodoHandler)

	err := http.ListenAndServe("localhost:8088", nil)
	utils.Check(err)
}
