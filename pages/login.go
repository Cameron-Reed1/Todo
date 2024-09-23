package pages

import (
	"net/http"

	"github.com/Cameron-Reed1/todo-web/auth"
	"github.com/Cameron-Reed1/todo-web/db"
	"github.com/Cameron-Reed1/todo-web/pages/templates"
)

func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        if _, err := validateSession(r); err == nil {
            w.Header().Add("Location", "/")
            w.WriteHeader(http.StatusSeeOther)
        } else {
            templates.LoginPage().Render(r.Context(), w)
        }
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")
    stay_logged := r.FormValue("stay-logged-in") == "on"

    if username == "" || password == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    userId, hashSalt, err := db.GetUserPassHash(username)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if auth.Validate(hashSalt.Hash, hashSalt.Salt, []byte(password)) {
        session, err := createSession(userId)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.Header().Add("Set-Cookie", session.ToCookie(stay_logged))
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusUnauthorized)
    }
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        if _, err := validateSession(r); err == nil {
            w.Header().Add("Location", "/")
            w.WriteHeader(http.StatusSeeOther)
        } else {
            templates.CreateAccountBox().Render(r.Context(), w)
        }
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    if username == "" || password == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // TODO: validate credentials
    // Ensure the username is valid and is not taken
    // Ensure that the password meets requirements

    hashSalt, err := auth.Hash([]byte(password), nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    user_id, err := db.CreateUser(username, hashSalt.Hash, hashSalt.Salt)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    session, err := createSession(user_id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Header().Add("Set-Cookie", session.ToCookie(false))
}

func Logout(w http.ResponseWriter, r *http.Request) {
    session, err := validateSession(r)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    err = db.DeleteSession(session)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Add("Set-Cookie", "session=;expires=Thu, 01 Jan 1970 00:00:00 UTC;samesite=strict;secure;HTTPonly")
    w.WriteHeader(http.StatusOK)
}
