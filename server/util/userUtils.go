package utils

import (
	"database/sql"

	"github.com/markbates/goth"
)

type User struct {
	ID        int
	Nickname  string
	Email     string
	AvatarURL string
}

var db = InitConn()

func InitUserTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY, 
		google_id VARCHAR(255) UNIQUE,
		email VARCHAR(255) UNIQUE NOT NULL,
		nickname VARCHAR(255),
		avatar_url TEXT,
		created timestamp DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	Check(err)
}

func SaveUserToDB(user goth.User) int {
	InitUserTable(db)
	var id int
	query := `INSERT INTO users (google_id, nickname, email, avatar_url)
		VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.QueryRow(query, user.UserID, user.NickName, user.Email, user.AvatarURL).Scan(&id)
	Check(err)
	return id
}

func GetUserFromDB(userid int) *User {
	user := &User{}
	err := db.QueryRow("SELECT id, nickname, email, avatar_url FROM users WHERE id = $1", userid).Scan(&user.ID, &user.Nickname, &user.Email, &user.AvatarURL)
	if err != nil {
		return nil
	}
	return user
}
