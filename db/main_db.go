package db

import (
	"database/sql"
	"encoding/hex"
	"log"

	"github.com/Cameron-Reed1/todo-web/auth"
	"github.com/Cameron-Reed1/todo-web/types"
	_ "github.com/mattn/go-sqlite3"
)


var main_db *sql.DB


func OpenMainDB(path string) {
    if main_db != nil {
        log.Fatal("Cannot open main DB twice!")
    }

    var err error
    main_db, err = sql.Open("sqlite3", path)
    if err != nil {
        log.Fatal(err)
    }

    err = main_db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    query := `
      CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        password_salt TEXT NOT NULL
      );
      CREATE TABLE IF NOT EXISTS sessions (
        sessionId TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
      );`

    _, err = main_db.Exec(query)
    if err != nil {
        log.Fatal(err)
    }
}

func CloseMainDB() {
    main_db.Close()
}

func CreateUser(username string, password_hash, password_salt []byte) (int64, error) {
    hex_hash := hex.EncodeToString(password_hash)
    hex_salt := hex.EncodeToString(password_salt)

    res, err := main_db.Exec("INSERT INTO users(username, password_hash, password_salt) values(?, ?, ?)", username, hex_hash, hex_salt)
    if err != nil {
        return 0, err
    }

    return res.LastInsertId()
}

func GetUserPassHash(username string) (int64, *auth.HashSalt, error) {
    hashSalt := auth.HashSalt{}
    var user_id int64
    var hex_hash string
    var hex_salt string

    row := main_db.QueryRow("SELECT id, password_hash, password_salt FROM users WHERE username=?", username)
    err := row.Scan(&user_id, &hex_hash, &hex_salt)
    if err != nil {
        return 0, nil, err
    }

    hashSalt.Hash, err = hex.DecodeString(hex_hash)
    if err != nil {
        return 0, nil, err
    }
    hashSalt.Salt, err = hex.DecodeString(hex_salt)
    if err != nil {
        return 0, nil, err
    }

    return user_id, &hashSalt, nil
}

func DeleteUser(username string) error {
    _, err := main_db.Exec("DELETE FROM users WHERE username=?", username)
    return err
}


func AddSession(session *types.Session) error {
    // fmt.Printf("New session: %s, %d\n", session.SessionId, session.UserId)
    _, err := main_db.Exec("INSERT INTO sessions(sessionId, user_id) values(?, ?)", session.SessionId, session.UserId)
    // fmt.Printf("Err: %v\n", err)
    return err
}

func GetUserFromSession(sessionId string) (string, error) {
    var username string

    row := main_db.QueryRow("SELECT username FROM sessions INNER JOIN users ON sessions.user_id = users.id WHERE sessionId=?", sessionId)
    err := row.Scan(&username)
    if err != nil {
        return "", err
    }

    return username, nil
}

func GetSession(sessionId string) (*types.Session, error) {
    var session types.Session

    row := main_db.QueryRow("SELECT user_id FROM sessions WHERE sessionId=?", sessionId)
    session.SessionId = sessionId
    err := row.Scan(&session.UserId)
    if err != nil {
        return nil, err
    }

    return &session, nil
}

func DeleteSession(sessionId string) error {
    _, err := main_db.Exec("DELETE FROM sessions WHERE sessionId=?", sessionId)
    return err
}
