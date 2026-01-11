package utils

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckLog(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}

type Note struct {
	Author string
	Text   string
}

type NotesData struct {
	NoteCount int
	Notes     []Note
}

func InitConn() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=gosha123X dbname=sitedb sslmode=disable"
	DB, err := sql.Open("postgres", connStr)
	Check(err)
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	return DB
}

func CloseConn() {
	if DB != nil {
		DB.Close()
	}
}

func CreateNoteTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY, 
		author VARCHAR(30) NOT NULL,
		text VARCHAR(120) NOT NULL,
		created timestamp DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	Check(err)
}

func InsertNoteTable(db *sql.DB, note Note) int {
	query := `INSERT INTO notes (author, text)
		VALUES ($1, $2) RETURNING id`

	var pk int
	err := db.QueryRow(query, note.Author, note.Text).Scan(&pk)
	Check(err)
	return pk
}

func AllNotesTable(db *sql.DB) []Note {
	data := []Note{}
	rows, err := db.Query("SELECT author, text FROM notes")
	Check(err)
	var author string
	var text string
	for rows.Next() {
		err := rows.Scan(&author, &text)
		Check(err)
		data = append(data, Note{author, text})
	}
	return data
}

type Task struct {
	ID   int
	Text string
	Done bool
}

type NewTask struct {
	Text string
	Done bool
}

type TasksData struct {
	TaskCount int
	Tasks     []Task
}

func CreateTaskTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY, 
		text VARCHAR(50) NOT NULL,
		done BOOL NOT NULL,
		created timestamp DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	Check(err)
}

func InsertTaskTable(db *sql.DB, task NewTask) int {
	query := `INSERT INTO tasks (text, done)
		VALUES ($1, $2) RETURNING id`

	var pk int
	err := db.QueryRow(query, task.Text, task.Done).Scan(&pk)
	Check(err)
	return pk
}

func UpdateTaskTable(db *sql.DB, id int, done bool) int {
	query := `UPDATE tasks SET done = $1 WHERE id = $2 RETURNING id`

	var pk int
	err := db.QueryRow(query, done, id).Scan(&pk)
	Check(err)
	return pk
}

func AllTaskTable(db *sql.DB) []Task {
	data := []Task{}
	rows, err := db.Query("SELECT id, text, done FROM tasks")
	Check(err)
	var id int
	var text string
	var done bool
	for rows.Next() {
		err := rows.Scan(&id, &text, &done)
		Check(err)
		data = append(data, Task{id, text, done})
	}
	return data
}
