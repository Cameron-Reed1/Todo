package pages

import (
	"net/http"

	"github.com/Cameron-Reed1/todo-web/pages/templates"
)

func OverdueFragment(w http.ResponseWriter, r *http.Request) {
    user_db, err := validateSessionAndGetUserDB(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    defer user_db.Close()

    items, err := user_db.GetOverdueTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList("Nothing to do in the past", items).Render(r.Context(), w)
}

func TodayFragment(w http.ResponseWriter, r *http.Request) {
    user_db, err := validateSessionAndGetUserDB(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    defer user_db.Close()

    items, err := user_db.GetTodayTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList("Nothing to do today", items).Render(r.Context(), w)
}

func UpcomingFragment(w http.ResponseWriter, r *http.Request) {
    user_db, err := validateSessionAndGetUserDB(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    defer user_db.Close()

    items, err := user_db.GetUpcomingTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList("Nothing to do in the future", items).Render(r.Context(), w)
}
