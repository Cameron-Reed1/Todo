package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Cameron-Reed1/todo-web/pages/templates"
	"github.com/Cameron-Reed1/todo-web/types"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
    user_db, err := validateSessionAndGetUserDB(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    defer user_db.Close()

    var todo types.Todo

    todo.Text = r.FormValue("name")
    start := r.FormValue("start")
    due := r.FormValue("due")

    // fmt.Printf("Create item request: %s: %s - %s\n", todo.Text, start, due)

    if start != "" {
        todo.Start, err = strconv.ParseInt(start, 10, 64)
        if err != nil {
            fmt.Printf("Bad start time: %s\n", start)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
    } else {
        todo.Start = 0
    }

    if due != "" {
        todo.Due, err = strconv.ParseInt(due, 10, 64)
        if err != nil {
            fmt.Printf("Bad due time: %s\n", due)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
    } else {
        todo.Due = 0
    }

    // fmt.Printf("New item: %s: %d - %d\n", todo.Text, todo.Start, todo.Due)

    err = user_db.AddTodo(&todo)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    now := time.Now().Unix()
    var targetSelector = "#today-list > .new-item"
    if todo.Due != 0 && todo.Due < now {
        targetSelector = "#overdue-list > .new-item"
    }
    if todo.Start > now {
        targetSelector = "#upcoming-list > .new-item"
    }

    templates.OobTodoItem(targetSelector, todo).Render(r.Context(), w)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
    user_db, err := validateSessionAndGetUserDB(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    defer user_db.Close()

    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)

    if idStr == "" || err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    err = user_db.DeleteTodo(id)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write([]byte{})
}

func SetItemCompleted(w http.ResponseWriter, r *http.Request) {
    user_db, err := validateSessionAndGetUserDB(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    defer user_db.Close()

    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)

    if idStr == "" || err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    completed := r.FormValue("completed") == "on"
    if err = user_db.SetCompleted(id, completed); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write([]byte{})
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
    user_db, err := validateSessionAndGetUserDB(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    defer user_db.Close()

    var todo types.Todo

    idStr := r.FormValue("id")
    todo.Text = r.FormValue("name")
    start := r.FormValue("start")
    due := r.FormValue("due")

    todo.Id, err = strconv.ParseInt(idStr, 10, 64)

    if idStr == "" || err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if start != "" {
        todo.Start, err = strconv.ParseInt(start, 10, 64)
        if err != nil {
            fmt.Printf("Bad start time: %s\n", start)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
    } else {
        todo.Start = 0
    }

    if due != "" {
        todo.Due, err = strconv.ParseInt(due, 10, 64)
        if err != nil {
            fmt.Printf("Bad due time: %s\n", due)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
    } else {
        todo.Due = 0
    }

    // fmt.Printf("New values:\n(%d) %s: %d - %d\n\n", todo.Id, todo.Text, todo.Start, todo.Due)

    err = user_db.UpdateTodo(todo)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    now := time.Now().Unix()
    var targetSelector = "#today-list > .new-item"
    if todo.Due != 0 && todo.Due < now {
        targetSelector = "#overdue-list > .new-item"
    }
    if todo.Start > now {
        targetSelector = "#upcoming-list > .new-item"
    }

    templates.OobTodoItem(targetSelector, todo).Render(r.Context(), w)
}
