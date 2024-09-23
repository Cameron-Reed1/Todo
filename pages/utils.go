package pages

import (
	"net/http"

	"github.com/Cameron-Reed1/todo-web/auth"
	"github.com/Cameron-Reed1/todo-web/db"
	"github.com/Cameron-Reed1/todo-web/types"
)

func createSession(user_id int64) (*types.Session, error) {
    session, err := auth.CreateSessionFor(user_id)
    if err != nil {
        return nil, err
    }

    err = db.AddSession(session)
    if err != nil {
        return nil, err
    }

    return session, nil
}

func validateSession(r *http.Request) (string, error) {
    cookie, err := r.Cookie("session")
    if err != nil {
        return "", err
    }

    _, err = db.GetUserFromSession(cookie.Value)
    // session, err := db.GetSession(cookie.Value)
    return cookie.Value, err
}

func validateSessionAndGetUsername(r *http.Request) (string, error) {
    cookie, err := r.Cookie("session")
    if err != nil {
        return "", err
    }

    return db.GetUserFromSession(cookie.Value)
}

func validateSessionAndGetUserDB(r *http.Request) (*db.UserDB, error) {
    cookie, err := r.Cookie("session")
    if err != nil {
        return nil, err
    }

    username, err := db.GetUserFromSession(cookie.Value)
    // session, err := db.GetSession(cookie.Value)
    if err != nil {
        return nil, err
    }

    user_db, err := db.OpenUserDB(username)
    if err != nil {
        return nil, err
    }

    return user_db, nil
}
