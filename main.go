package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/Cameron-Reed1/todo-web/api"
	"github.com/Cameron-Reed1/todo-web/db"
	"github.com/Cameron-Reed1/todo-web/pages"
)

func main() {
    db_path := flag.String("db", "./test.db", "Path to the sqlite3 database")
    bind_port := flag.Int("p", 8080, "Port to bind to")
    bind_addr := flag.String("a", "0.0.0.0", "Address to bind to")
    noFront := flag.Bool("no-frontend", false, "Disable the frontend endpoints")
    a := false; noBack := &a // flag.Bool("no-backend", false, "Disable the backend endpoints")

    flag.Parse()

    mux := http.NewServeMux()

    if *noFront && *noBack {
        fmt.Println("What do you want me to do?")
        return
    }

    if !*noFront {
        addFrontendEndpoints(mux)
    }

    if !*noBack {
        addBackendEndpoints(mux)
    }

    db.Open(*db_path)
    defer db.Close()

    addr := fmt.Sprintf("%s:%d", *bind_addr, *bind_port)
    server := http.Server{ Addr: addr, Handler: mux }
    fmt.Printf("Starting server on %s...\n", addr)
    err := server.ListenAndServe()

    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("Server closed\n")
    } else if err != nil {
        fmt.Printf("Error starting server: %s\n", err)
    }
}

func addFrontendEndpoints(mux *http.ServeMux) {
    fmt.Println("Frontend enabled")

    mux.HandleFunc("/", Error404)

    mux.HandleFunc("/{$}", pages.RootPage)
    mux.HandleFunc("/overdue", pages.OverdueFragment)
    mux.HandleFunc("/today", pages.TodayFragment)
    mux.HandleFunc("/upcoming", pages.UpcomingFragment)
    mux.HandleFunc("DELETE /delete/{id}", pages.DeleteItem)
    mux.HandleFunc("PATCH /set/{id}", pages.SetItemCompleted)
    mux.HandleFunc("POST /new", pages.CreateItem)

    fileServer := http.FileServer(http.Dir("./static"))
    mux.Handle("/css/", fileServer)
    mux.Handle("/js/", fileServer)
}

func addBackendEndpoints(mux *http.ServeMux) {
    fmt.Println("Backend enabled")

    mux.HandleFunc("/api/", api.InvalidEndpoint)

    mux.HandleFunc("GET /api/get", api.GetAll)
    mux.HandleFunc("GET /api/get/{id}", api.GetTodo)
    mux.HandleFunc("POST /api/new", api.AddTodo)
}

func Error404(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Error 404: Page not found\n"))
}
