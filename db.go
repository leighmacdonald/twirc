package twirc

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)

var (
	SqlDB  *sqlx.DB
	tables = map[string]string{
		"user": `CREATE TABLE user (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			steamid TEXT UNIQUE,
			created_at DATETIME
		)`,
		"messages": `CREATE TABLE messages (
			message_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			message TEXT NOT NULL,
			emotes INTEGER,
			created_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES user(user_id)
		)`,
	}
)

type ChatMsg struct {
	MessageID int       `db:message_id`
	UserID    int       `db:user_id`
	Message   string    `db:message`
	Emotes    int       `db:emotes`
	CreatedAt time.Time `db:created_at`
}

type User struct {
	UserID    int       `db:user_id`
	Username  string    `db:username`
	SteamID   string    `db:steamid`
	CreatedAt time.Time `db:created_at`
}

func NewChatMsg(user *User, message string) *ChatMsg {
	return &ChatMsg{
		UserID:    user.UserID,
		Message:   message,
		Emotes:    1,
		CreatedAt: time.Now(),
	}
}

func (msg *ChatMsg) Save(db *sqlx.DB) {
	q := "INSERT INTO messages (user_id, message, emotes, created_at) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(q, msg.UserID, msg.Message, msg.Emotes, time.Now())
	if err != nil {
		log.Println(err.Error())
	}
}

func initTables(d *sqlx.DB) {
	for table_name, create_stmt := range tables {
		stmt, err := d.Prepare("SELECT name FROM sqlite_master WHERE type='table' AND name=?;")
		if err != nil {
			panic(err)
		}
		defer stmt.Close()
		var name string
		err = stmt.QueryRow(table_name).Scan(&name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				log.Printf("Initializing table: %s\n", table_name)
				if _, err := d.Exec(create_stmt); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
	}
}

func GetUser(username string) (*User, error) {
	var user User
	row := SqlDB.QueryRow("SELECT user_id, username, steamid, created_at FROM user WHERE username = ? LIMIT 1", strings.ToLower(username))
	err := row.Scan(&user.UserID, &user.Username, &user.SteamID, &user.CreatedAt)
	PrettyPrint(user)
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
	return &user, err
}

func GetOrCreateUser(username string) *User {
	user, err := GetUser(username)
	if err != nil {
		q := "INSERT INTO user (username, created_at) VALUES (?, ?)"
		_, err := SqlDB.Exec(q, strings.ToLower(username), time.Now())
		if err != nil {
			log.Println(err.Error())
		}
	}
	return user
}

func DeleteUserByName(db *sqlx.DB, username string) bool {
	rows, err := db.MustExec("DELETE FROM user WHERE username = ?", username).RowsAffected()
	if err != nil {
		log.Printf(err.Error())
		return false
	}
	if rows >= 1 {
		log.Printf("Deleted user: %s\n", username)
	}
	return rows >= 1
}
