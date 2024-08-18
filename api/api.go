package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Cameron-Reed1/todo-web/db"
	"github.com/Cameron-Reed1/todo-web/types"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
    todos, err := db.GetAllTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("{\"error\":\"Failed to get items\"}"))
        return
    }

    response := "{\"items\":["

    first := true
    for _, todo := range todos {
        if !first {
            response += ","
        }
        str := fmt.Sprintf("{\"id\":%d,\"start\":%d,\"due\":%d,\"text\":\"%s\"}", todo.Id, todo.Start, todo.Due, todo.Text)
        response += str
        first = false
    }

    response += "]}"
    w.Write([]byte(response))
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)

    if idStr == "" || err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Invalid id\"}"))
        return
    }

    todo, err := db.GetTodo(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("{\"error\":\"No item for id\"}"))
        return
    }

    str := fmt.Sprintf("{\"id\":%d,\"start\":%d,\"due\":%d,\"text\":\"%s\"}", todo.Id, todo.Start, todo.Due, todo.Text)
    w.Write([]byte(str))
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
    var todo types.Todo

    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    err := decoder.Decode(&todo)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Failed to parse JSON\"}"))
        return
    }

    if decoder.More() {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Extra data after JSON object\"}"))
        return
    }

    if todo.Text == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Invalid text\"}"))
        return
    }

    err = db.AddTodo(&todo)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("{\"error\":\"Failed to add item\"}"))
        return
    }

    res := fmt.Sprintf("{\"id\":%d}", todo.Id)
    w.Write([]byte(res))
}

func InvalidEndpoint(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("{\"error\":\"Endpoint not found\"}"))
}
