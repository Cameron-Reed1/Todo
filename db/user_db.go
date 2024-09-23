package db

import (
	"database/sql"
	"errors"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Cameron-Reed1/todo-web/types"
)

var userDbDir string

type UserDB struct {
    DB *sql.DB
}

func SetUserDBDir(dir string) error {
    os.MkdirAll(dir, 0700)
    userDbDir = dir
    return nil
}

func OpenUserDB(username string) (*UserDB, error) {
    if strings.Contains(username, ".") || strings.Contains(username, "/") {
        return nil, errors.New("Invalid username")
    }

    path := path.Join(userDbDir, username + ".db")

    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    query := `
      CREATE TABLE IF NOT EXISTS items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        start INTEGER,
        due INTEGER,
        text TEXT NOT NULL,
        completed INTEGER NOT NULL
      );`

    _, err = db.Exec(query)

    return &UserDB{DB: db}, err
}

func (user_db *UserDB) Close() error {
    return user_db.DB.Close()
}

func (user_db *UserDB) AddTodo(todo *types.Todo) error {
    res, err := user_db.DB.Exec("INSERT INTO items(start, due, text, completed) values(?, ?, ?, 0)", toNullInt64(todo.Start), toNullInt64(todo.Due), todo.Text)
    if err != nil {
        return err
    }

    todo.Id, err = res.LastInsertId()
    if err != nil {
        return err
    }

    return nil
}

func (user_db *UserDB) GetTodo(id int) (types.Todo, error) {
    var todo types.Todo
    var start sql.NullInt64
    var due sql.NullInt64

    row := user_db.DB.QueryRow("SELECT * FROM items WHERE id=?", id)
    err := row.Scan(&todo.Id, &start, &due, &todo.Text, &todo.Completed)

    todo.Start = fromNullInt64(start)
    todo.Due = fromNullInt64(due)

    return todo, err
}

func (user_db *UserDB) GetAllTodos() ([]types.Todo, error) {
    var todos []types.Todo

    rows, err := user_db.DB.Query("SELECT * FROM items")
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var todo types.Todo
        var start sql.NullInt64
        var due sql.NullInt64

        err = rows.Scan(&todo.Id, &start, &due, &todo.Text, &todo.Completed)
        if err != nil {
            return nil, err
        }

        todo.Start = fromNullInt64(start)
        todo.Due = fromNullInt64(due)

        todos = append(todos, todo)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}

func (user_db *UserDB) GetOverdueTodos() ([]types.Todo, error) {
    var todos []types.Todo

    rows, err := user_db.DB.Query("SELECT * FROM items WHERE due < ? AND due IS NOT NULL ORDER BY completed, due", time.Now().Unix())
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var todo types.Todo
        var start sql.NullInt64
        var due sql.NullInt64

        err = rows.Scan(&todo.Id, &start, &due, &todo.Text, &todo.Completed)
        if err != nil {
            return nil, err
        }

        todo.Start = fromNullInt64(start)
        todo.Due = fromNullInt64(due)

        todos = append(todos, todo)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}

func (user_db *UserDB) GetTodayTodos() ([]types.Todo, error) {
    var todos []types.Todo

    now := time.Now()
    year, month, day := now.Date()
    today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
    rows, err := user_db.DB.Query("SELECT * FROM items WHERE (start <= ? OR start IS NULL) AND (due >= ? OR due IS NULL) ORDER BY completed, due NULLS LAST", today.Unix(), now.Unix())
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var todo types.Todo
        var start sql.NullInt64
        var due sql.NullInt64

        err = rows.Scan(&todo.Id, &start, &due, &todo.Text, &todo.Completed)
        if err != nil {
            return nil, err
        }

        todo.Start = fromNullInt64(start)
        todo.Due = fromNullInt64(due)

        todos = append(todos, todo)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}

func (user_db *UserDB) GetUpcomingTodos() ([]types.Todo, error) {
    var todos []types.Todo

    year, month, day := time.Now().Date()
    today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
    rows, err := user_db.DB.Query("SELECT * FROM items WHERE start > ? ORDER BY completed, start", today.Unix())
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var todo types.Todo
        var start sql.NullInt64
        var due sql.NullInt64

        err = rows.Scan(&todo.Id, &start, &due, &todo.Text, &todo.Completed)
        if err != nil {
            return nil, err
        }

        todo.Start = fromNullInt64(start)
        todo.Due = fromNullInt64(due)

        todos = append(todos, todo)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}

func (user_db *UserDB) UpdateTodo(newValues types.Todo) error {
    _, err := user_db.DB.Exec("UPDATE items SET start=?, due=?, text=? WHERE id=?", toNullInt64(newValues.Start), toNullInt64(newValues.Due), newValues.Text, newValues.Id)
    return err;
}

func (user_db *UserDB) SetCompleted(id int, completed bool) error {
    _, err := user_db.DB.Exec("UPDATE items SET completed=? WHERE id=?", completed, id)
    return err
}

func (user_db *UserDB) DeleteTodo(id int) error {
    _, err := user_db.DB.Exec("DELETE FROM items WHERE id=?", id)
    return err
}
