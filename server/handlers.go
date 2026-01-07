package handlers

import (
	"html/template"
	"net/http"
	utils "site/server/util"
	"strconv"
)

func MainHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("static/templates/main.html")
	utils.Check(err)
	err = html.Execute(writer, nil)
	utils.Check(err)
}

func NotesHandler(writer http.ResponseWriter, request *http.Request) {
	db := utils.InitConn()
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
	utils.CloseConn()
}

func CreateNoteHandler(writer http.ResponseWriter, request *http.Request) {
	db := utils.InitConn()
	textInput := request.FormValue("textInput")
	authorInput := request.FormValue("authorInput")
	newNote := utils.Note{Author: authorInput, Text: textInput}
	utils.InsertNoteTable(db, newNote)
	http.Redirect(writer, request, "/notes", http.StatusFound)
	utils.CloseConn()
}

func TodoHandler(writer http.ResponseWriter, request *http.Request) {
	db := utils.InitConn()
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
	utils.CloseConn()
}

func CreateTodoHandler(writer http.ResponseWriter, request *http.Request) {
	db := utils.InitConn()
	taskInput := request.FormValue("taskInput")
	newTask := utils.NewTask{Text: taskInput, Done: false}
	utils.InsertTaskTable(db, newTask)
	http.Redirect(writer, request, "/todo", http.StatusFound)
	utils.CloseConn()
}

func UpdateTodoHandler(writer http.ResponseWriter, request *http.Request) {
	db := utils.InitConn()
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
	utils.CloseConn()
}
