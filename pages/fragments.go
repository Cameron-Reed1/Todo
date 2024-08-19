package pages

import (
	"net/http"

	"github.com/Cameron-Reed1/todo-web/db"
	"github.com/Cameron-Reed1/todo-web/pages/templates"
)

func OverdueFragment(w http.ResponseWriter, r *http.Request) {
    items, err := db.GetOverdueTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList("Nothing to do in the past", items).Render(r.Context(), w)
}

func TodayFragment(w http.ResponseWriter, r *http.Request) {
    items, err := db.GetTodayTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList("Nothing to do today", items).Render(r.Context(), w)
}

func UpcomingFragment(w http.ResponseWriter, r *http.Request) {
    items, err := db.GetUpcomingTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList("Nothing to do in the future", items).Render(r.Context(), w)
}
