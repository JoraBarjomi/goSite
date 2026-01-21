package main

import (
	"log"
	"net/http"
	"os"
	handlers "site/server"
	utils "site/server/util"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	initEnv()
	var key = os.Getenv("KEY")
	var store = sessions.NewCookieStore([]byte(key))

	store.Options.HttpOnly = true
	store.Options.Secure = true

	gothic.Store = store

	handlers.InitGoogle()
	http.HandleFunc("/auth/google", handlers.GoogleAuthHandler)
	http.HandleFunc("/auth/google/callback", handlers.GoogleCallbackHandler)

	http.HandleFunc("/profile", handlers.ProfileHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/login/auth", handlers.LoginHelperHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/register/create", handlers.CreateUserHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	http.HandleFunc("/", handlers.MainHandler)

	http.HandleFunc("/notes", handlers.NotesHandler)
	http.HandleFunc("/notes/create", handlers.CreateNoteHandler)

	http.HandleFunc("/todo", handlers.TodoHandler)
	http.HandleFunc("/todo/create", handlers.CreateTodoHandler)
	http.HandleFunc("/todo/update", handlers.UpdateTodoHandler)

	hub := newHub()
	go hub.run()
	http.HandleFunc("/chat", handlers.WsHelperHandler)
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(hub, writer, request)
	})

	defer utils.CloseConn()

	err := http.ListenAndServe("localhost:8088", nil)
	utils.Check(err)
}
