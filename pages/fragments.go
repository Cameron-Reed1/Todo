package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Cameron-Reed1/todo-web/db"
	"github.com/Cameron-Reed1/todo-web/pages/templates"
	"github.com/Cameron-Reed1/todo-web/types"
)

func OverdueFragment(w http.ResponseWriter, r *http.Request) {
    items, err := db.GetOverdueTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList(items).Render(r.Context(), w)
}

func TodayFragment(w http.ResponseWriter, r *http.Request) {
    items, err := db.GetTodayTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList(items).Render(r.Context(), w)
}

func UpcomingFragment(w http.ResponseWriter, r *http.Request) {
    items, err := db.GetUpcomingTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    templates.TodoList(items).Render(r.Context(), w)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
    var todo types.Todo
    var err error

    todo.Text = r.FormValue("name")
    start := r.FormValue("start")
    due := r.FormValue("due")

    if start != "" {
        start_time, err := time.Parse("2006-01-02T03:04", start)
        if err != nil {
            fmt.Printf("Bad start time: %s\n", start)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        todo.Start = start_time.Unix()
    } else {
        todo.Start = 0
    }

    if due != "" {
        due_time, err := time.Parse("2006-01-02T15:04", due)
        if err != nil {
            fmt.Printf("Bad due time: %s\n", due)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        todo.Due = due_time.Unix()
    } else {
        todo.Due = 0
    }

    fmt.Printf("New item: %s: %d - %d\n", todo.Text, todo.Start, todo.Due)

    err = db.AddTodo(&todo)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write([]byte{})
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)

    if idStr == "" || err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    err = db.DeleteTodo(id)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write([]byte{})
}

func SetItemCompleted(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)

    if idStr == "" || err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    completed := r.FormValue("completed") == "on"
    if err = db.SetCompleted(id, completed); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write([]byte{})
}
