package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	utils "site/server/util"
)

func MainHandler(writer http.ResponseWriter, request *http.Request) {
	notes := utils.GetStrings("static/database/blogs.txt")
	html, err := template.ParseFiles("static/templates/main.html")
	utils.Check(err)
	note := utils.Note{
		NoteCount: len(notes),
		Notes:     notes,
	}
	err = html.Execute(writer, note)
	utils.Check(err)
}

func NewHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("static/templates/new.html")
	utils.Check(err)
	err = html.Execute(writer, nil)
	utils.Check(err)
}

func CreateHandler(writer http.ResponseWriter, request *http.Request) {
	blogInput := request.FormValue("inputBlog")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("static/database/blogs.txt", options, os.FileMode(0600))
	utils.Check(err)
	_, err = fmt.Fprintln(file, blogInput)
	utils.Check(err)
	err = file.Close()
	utils.Check(err)
	http.Redirect(writer, request, "/notes", http.StatusFound)
}
