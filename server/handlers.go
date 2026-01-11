package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	utils "site/server/util"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var db = utils.InitConn()

func MainHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("static/templates/main.html")
	utils.Check(err)
	session, _ := gothic.Store.Get(request, "app_session")
	userID, ok := session.Values["user_id"]
	if !ok {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	user := utils.GetUserFromDB(userID.(int))
	if user == nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	log.Println("Hello, ", user.Nickname)

	err = html.Execute(writer, user)
	utils.Check(err)
}

var upgrader = websocket.Upgrader{}

func WsHandler(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	utils.CheckLog(err)
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(mt, message)
	}
}

func WsHelperHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "static/templates/ppchat.html")
}

func NotesHandler(writer http.ResponseWriter, request *http.Request) {
	utils.CreateNoteTable(db)
	notes := utils.AllNotesTable(db)
	html, err := template.ParseFiles("static/templates/notes.html")
	utils.Check(err)
	note := utils.NotesData{
		NoteCount: len(notes),
		Notes:     notes,
	}
	err = html.Execute(writer, note)
	utils.Check(err)
}

func CreateNoteHandler(writer http.ResponseWriter, request *http.Request) {
	textInput := request.FormValue("textInput")
	authorInput := request.FormValue("authorInput")
	newNote := utils.Note{Author: authorInput, Text: textInput}
	utils.InsertNoteTable(db, newNote)
	http.Redirect(writer, request, "/notes", http.StatusFound)
}

func TodoHandler(writer http.ResponseWriter, request *http.Request) {
	utils.CreateTaskTable(db)
	tasks := utils.AllTaskTable(db)
	html, err := template.ParseFiles("static/templates/todo.html")
	utils.Check(err)
	task := utils.TasksData{
		TaskCount: len(tasks),
		Tasks:     tasks,
	}
	err = html.Execute(writer, task)
	utils.Check(err)
}

func CreateTodoHandler(writer http.ResponseWriter, request *http.Request) {
	taskInput := request.FormValue("taskInput")
	newTask := utils.NewTask{Text: taskInput, Done: false}
	utils.InsertTaskTable(db, newTask)
	http.Redirect(writer, request, "/todo", http.StatusFound)
}

func UpdateTodoHandler(writer http.ResponseWriter, request *http.Request) {
	idTask := request.FormValue("id")
	id, err := strconv.Atoi(idTask)
	utils.Check(err)
	doneTask := request.FormValue("done")
	var done bool
	if doneTask == "on" {
		done = true
	} else {
		done = false
	}
	utils.UpdateTaskTable(db, id, done)
	http.Redirect(writer, request, "/todo", http.StatusFound)
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("static/templates/auth/login.html")
	utils.Check(err)
	err = html.Execute(writer, nil)
	utils.Check(err)
}

func LogoutHandler(writer http.ResponseWriter, request *http.Request) {
	session, _ := gothic.Store.Get(request, "app_session")
	session.Options.MaxAge = -1
	session.Values = make(map[interface{}]interface{})
	session.Save(request, writer)

	log.Println("User logged out")
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}

func InitGoogle() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientId == "" && googleClientSecret == "" {
		log.Fatal("Error parsing google variables")
	}
	goth.UseProviders(google.New(googleClientId, googleClientSecret, "http://localhost:8088/auth/google/callback"))
}

func GoogleAuthHandler(writer http.ResponseWriter, request *http.Request) {
	session, _ := gothic.Store.Get(request, "app_session")
	if _, ok := session.Values["user_id"]; ok {
		log.Println("User already logged in:")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	q := request.URL.Query()
	q.Add("provider", "google")
	request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(writer, request)
}

func GoogleCallbackHandler(writer http.ResponseWriter, request *http.Request) {
	user, err := gothic.CompleteUserAuth(writer, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}
	log.Println("New user: ", user.NickName)
	dbUserId := utils.SaveUserToDB(user)
	session, _ := gothic.Store.Get(request, "app_session")
	session.Values["user_id"] = dbUserId
	session.Save(request, writer)

	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
