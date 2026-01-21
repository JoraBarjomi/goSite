package utils

import (
	"database/sql"
	"strings"

	"github.com/markbates/goth"
)

type User struct {
	ID        int
	Nickname  string
	Email     string
	AvatarURL string
}

type NewUser struct {
	Id        int
	Nickname  string
	Email     string
	AvatarEnc string
	Password  string
}

var db = InitConn()

func InitUserTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY, 
		google_id VARCHAR(255) UNIQUE,
		email VARCHAR(255) UNIQUE NOT NULL,
		nickname VARCHAR(255),
		avatar_url TEXT,
		password VARCHAR(255),
		created timestamp DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	Check(err)
}

func SaveUserToDB(user goth.User) int { //for google auth
	InitUserTable(db)
	var id int
	query := `
		INSERT INTO users (google_id, nickname, email, avatar_url)
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (google_id)
		DO UPDATE SET
			nickname = EXCLUDED.nickname,
			email = EXCLUDED.email,
			avatar_url = EXCLUDED.avatar_url
		RETURNING id`
	err := db.QueryRow(query, user.UserID, user.NickName, user.Email, user.AvatarURL).Scan(&id)
	Check(err)
	return id
}

func SaveUserToDBReg(user NewUser) int {
	InitUserTable(db)
	var id int
	query := `
		INSERT INTO users (nickname, email, avatar_url, password)
		VALUES ($1, $2, $3, $4) 
		RETURNING id`
	err := db.QueryRow(query, user.Nickname, user.Email, user.AvatarEnc, user.Password).Scan(&id)
	Check(err)
	user.Id = id
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

func GetUserAvatar(user *User) string {
	if user == nil || user.AvatarURL == "" {
		return ""
	}

	if strings.HasPrefix(user.AvatarURL, "http") {
		return user.AvatarURL
	}

	return "data:image/jpeg;base64," + user.AvatarURL
}

func GetUserByEmail(email string) *NewUser {
	user := &NewUser{}
	err := db.QueryRow("SELECT id, nickname, email, avatar_url, password FROM users WHERE email = $1", email).Scan(&user.Id, &user.Nickname, &user.Email, &user.AvatarEnc, &user.Password)
	if err != nil {
		return nil
	}
	return user
}
