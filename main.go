package main

import (
	"net/http"
	handlers "site/server"
	utils "site/server/util"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

var store = sessions.NewCookieStore([]byte("_VLkokfwpK8DlARCsBIpINbZQTxRgGkThHC0D8Dtik0"))

func main() {

	gothic.Store = store

	handlers.InitGoogle()
	http.HandleFunc("/auth/google", handlers.GoogleAuthHandler)
	http.HandleFunc("/auth/google/callback", handlers.GoogleCallbackHandler)

	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	http.HandleFunc("/", handlers.MainHandler)

	http.HandleFunc("/notes", handlers.NotesHandler)
	http.HandleFunc("/notes/create", handlers.CreateNoteHandler)

	http.HandleFunc("/todo", handlers.TodoHandler)
	http.HandleFunc("/todo/create", handlers.CreateTodoHandler)
	http.HandleFunc("/todo/update", handlers.UpdateTodoHandler)

	http.Handle("/chat", http.HandlerFunc(handlers.WsHelperHandler))
	http.HandleFunc("/ws", handlers.WsHandler)

	defer utils.CloseConn()

	err := http.ListenAndServe("localhost:8088", nil)
	utils.Check(err)
}
