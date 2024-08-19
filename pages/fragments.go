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

    templates.TodoList("No overdue items", items).Render(r.Context(), w)
}

func TodayFragment(w http.ResponseWriter, r *http.Request) {
    items, err := db.GetTodayTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList("No items for today", items).Render(r.Context(), w)
}

func UpcomingFragment(w http.ResponseWriter, r *http.Request) {
    items, err := db.GetUpcomingTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList("No upcoming items", items).Render(r.Context(), w)
}
