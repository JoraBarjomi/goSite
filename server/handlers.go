package handlers

import (
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	utils "site/server/util"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var db = utils.InitConn()

type Base struct {
	User       *utils.User
	IsLoggedIn bool
	Title      string
	AvatarSrc  string
}

type PageNote struct {
	Base
	NoteCount int
	Notes     []utils.Note
}

type PageTodo struct {
	Base
	TaskCount int
	Tasks     []utils.Task
}

//Main page

func MainHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := gothic.Store.Get(request, "app_session")
	utils.Check(err)
	userID, ok := session.Values["user_id"]

	var user *utils.User
	isLoggedIn := false

	if ok {
		user = utils.GetUserFromDB(userID.(int))
		if user != nil {
			isLoggedIn = true
			log.Printf("User %s in db and logged in!", user.Email)
		}
	}

	data := Base{
		User:       user,
		IsLoggedIn: isLoggedIn,
		Title:      "Main - Site",
	}

	tmpl := template.Must(template.ParseFiles("static/templates/base.html", "static/templates/main.html"))
	err = tmpl.ExecuteTemplate(writer, "base.html", data)
	utils.Check(err)
}

//Ws chat page

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
	session, err := gothic.Store.Get(request, "app_session")
	utils.Check(err)
	userID, ok := session.Values["user_id"]

	var user *utils.User
	isLoggedIn := false

	if ok {
		user = utils.GetUserFromDB(userID.(int))
		if user != nil {
			isLoggedIn = true
			log.Printf("User %s in the chat!", user.Email)
		}
	}

	data := Base{
		User:       user,
		IsLoggedIn: isLoggedIn,
		Title:      "Chat - Site",
	}

	tmpl := template.Must(template.ParseFiles("static/templates/baseChat.html", "static/templates/ppchat.html"))
	err = tmpl.ExecuteTemplate(writer, "baseChat.html", data)
	utils.Check(err)
}

//Note page

func NotesHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := gothic.Store.Get(request, "app_session")
	utils.Check(err)
	userID, ok := session.Values["user_id"]

	var user *utils.User
	isLoggedIn := false

	if ok {
		user = utils.GetUserFromDB(userID.(int))
		if user != nil {
			isLoggedIn = true
			log.Printf("User %s in the chat!", user.Email)
		}
	}

	utils.CreateNoteTable(db)
	notes := utils.AllNotesTable(db)

	data := PageNote{
		Base: Base{
			User:       user,
			IsLoggedIn: isLoggedIn,
			Title:      "Chat - Site",
		},
		NoteCount: len(notes),
		Notes:     notes,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/base.html", "static/templates/notes.html"))
	err = tmpl.ExecuteTemplate(writer, "base.html", data)
	utils.Check(err)
}

func CreateNoteHandler(writer http.ResponseWriter, request *http.Request) {
	textInput := request.FormValue("textInput")
	authorInput := request.FormValue("authorInput")
	newNote := utils.Note{Author: authorInput, Text: textInput}
	utils.InsertNoteTable(db, newNote)
	http.Redirect(writer, request, "/notes", http.StatusFound)
}

// Todo page

func TodoHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := gothic.Store.Get(request, "app_session")
	utils.Check(err)
	userID, ok := session.Values["user_id"]

	var user *utils.User
	isLoggedIn := false

	if ok {
		user = utils.GetUserFromDB(userID.(int))
		if user != nil {
			isLoggedIn = true
			log.Printf("User %s in the chat!", user.Email)
		}
	}

	utils.CreateNoteTable(db)
	tasks := utils.AllTaskTable(db)

	data := PageTodo{
		Base: Base{
			User:       user,
			IsLoggedIn: isLoggedIn,
			Title:      "Chat - Site",
		},
		TaskCount: len(tasks),
		Tasks:     tasks,
	}
	tmpl := template.Must(template.ParseFiles("static/templates/base.html", "static/templates/todo.html"))
	err = tmpl.ExecuteTemplate(writer, "base.html", data)
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

//Auth handlers

func RegisterHandler(writer http.ResponseWriter, request *http.Request) {

	tmpl := template.Must(template.ParseFiles("static/templates/base.html", "static/templates/auth/register.html"))
	err := tmpl.ExecuteTemplate(writer, "base.html", nil)
	utils.Check(err)

}

func CreateUserHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20)
	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")
	passwordConf := request.FormValue("passwordConf")

	log.Println("PAROLI: ", password, "PAROL 2: ", passwordConf)
	log.Printf("Method: %s", request.Method)
	log.Printf("Content-Type: %s", request.Header.Get("Content-Type"))

	if password != passwordConf {
		http.Error(writer, "Error: password mismatch!", http.StatusBadRequest)
		return
	}

	if !strings.ContainsAny(password, "!@#$%^&*") {
		http.Error(writer, "Error: password not contain special symbol!", http.StatusBadRequest)
		return
	}

	file, _, err := request.FormFile("avatar")
	utils.Check(err)
	defer file.Close()
	data, _ := io.ReadAll(file)
	avatar := base64.StdEncoding.EncodeToString(data)

	user := utils.NewUser{Nickname: username, Email: email, AvatarEnc: avatar, Password: password}

	id := utils.SaveUserToDBReg(user)
	session, _ := gothic.Store.Get(request, "app_session")
	session.Values["user_id"] = id
	session.Save(request, writer)

	http.Redirect(writer, request, "/profile", http.StatusSeeOther)
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := gothic.Store.Get(request, "app_session")
	utils.Check(err)
	userID, ok := session.Values["user_id"]

	var user *utils.User
	isLoggedIn := false

	if ok {
		user = utils.GetUserFromDB(userID.(int))
		if user != nil {
			isLoggedIn = true
			log.Printf("User %s in login page!", user.Email)
		}
	}

	data := Base{
		User:       user,
		IsLoggedIn: isLoggedIn,
		Title:      "Login - Site",
	}

	tmpl := template.Must(template.ParseFiles("static/templates/base.html", "static/templates/auth/login.html"))
	err = tmpl.ExecuteTemplate(writer, "base.html", data)
	utils.Check(err)
}

func LoginHelperHandler(writer http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")
	password := request.FormValue("password")

	log.Println("EMAIL : ", email, "PSWD: ", password)

	user := utils.GetUserByEmail(email)
	if user == nil {
		log.Println("User not found!")
		http.Error(writer, "User not found!", http.StatusUnauthorized)
		return
	}
	log.Printf("Found user in db: %v %v %v %v ", user.Id, user.Nickname, user.Email, user.Password)

	log.Println(user.Email, user.Password)
	if user.Password != password && password != "" {
		http.Error(writer, "Not login!", http.StatusUnauthorized)
		return
	}

	session, _ := gothic.Store.Get(request, "app_session")
	session.Values["user_id"] = user.Id
	session.Save(request, writer)
	log.Println("Session is saved id: !", user.Id)

	http.Redirect(writer, request, "/profile", http.StatusSeeOther)
}

func LogoutHandler(writer http.ResponseWriter, request *http.Request) {
	session, _ := gothic.Store.Get(request, "app_session")
	session.Options.MaxAge = -1
	session.Values = make(map[interface{}]interface{})
	session.Save(request, writer)

	log.Println("User logged out")
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}

func ProfileHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := gothic.Store.Get(request, "app_session")
	utils.Check(err)
	userID, ok := session.Values["user_id"]

	if !ok || userID == nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	user := utils.GetUserFromDB(userID.(int))
	if user == nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	data := Base{
		User:       user,
		IsLoggedIn: true,
		Title:      "Profile - Site",
		AvatarSrc:  utils.GetUserAvatar(user),
	}

	tmpl := template.Must(template.ParseFiles("static/templates/base.html", "static/templates/auth/profile.html"))
	err = tmpl.ExecuteTemplate(writer, "base.html", data)
	utils.Check(err)
}

//Google auth

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
	log.Println("New user: ", user.Email)
	dbUserId := utils.SaveUserToDB(user)
	session, _ := gothic.Store.Get(request, "app_session")
	session.Values["user_id"] = dbUserId
	session.Save(request, writer)

	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
