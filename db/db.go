package db

import (
	"database/sql"
	"log"
    "time"

	_ "github.com/mattn/go-sqlite3"

    "github.com/Cameron-Reed1/todo-web/types"
)

var db *sql.DB

func Open(path string) {
    if db != nil {
        log.Fatal("Cannot init DB twice!")
    }

    var err error
    db, err = sql.Open("sqlite3", path)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
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
}

func AddTodo(todo *types.Todo) error {
    res, err := db.Exec("INSERT INTO items(start, due, text, completed) values(?, ?, ?, 0)", toNullInt64(todo.Start), toNullInt64(todo.Due), todo.Text)
    if err != nil {
        return err
    }

    todo.Id, err = res.LastInsertId()
    if err != nil {
        return err
    }

    return nil
}

func GetTodo(id int) (types.Todo, error) {
    var todo types.Todo
    var start sql.NullInt64
    var due sql.NullInt64

    row := db.QueryRow("SELECT * FROM items WHERE id=?", id)
    err := row.Scan(&todo.Id, &start, &due, &todo.Text, &todo.Completed)

    todo.Start = fromNullInt64(start)
    todo.Due = fromNullInt64(due)

    return todo, err
}

func GetAllTodos() ([]types.Todo, error) {
    var todos []types.Todo

    rows, err := db.Query("SELECT * FROM items")
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

func GetOverdueTodos() ([]types.Todo, error) {
    var todos []types.Todo

    rows, err := db.Query("SELECT * FROM items WHERE due < ? AND due IS NOT NULL ORDER BY due", time.Now().Unix())
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

func GetTodayTodos() ([]types.Todo, error) {
    var todos []types.Todo

    now := time.Now()
    year, month, day := now.Date()
    today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
    rows, err := db.Query("SELECT * FROM items WHERE (start <= ? OR start IS NULL) AND (due >= ? OR due IS NULL) ORDER BY due NULLS LAST", today.Unix(), now.Unix())
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

func GetUpcomingTodos() ([]types.Todo, error) {
    var todos []types.Todo

    year, month, day := time.Now().Date()
    today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
    rows, err := db.Query("SELECT * FROM items WHERE start > ? ORDER BY start", today.Unix())
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

func UpdateTodo(newValues types.Todo) error {
    _, err := db.Exec("UPDATE items SET start=?, due=?, text=? WHERE id=?", toNullInt64(newValues.Start), toNullInt64(newValues.Due), newValues.Text, newValues.Id)
    return err;
}

func SetCompleted(id int, completed bool) error {
    _, err := db.Exec("UPDATE items SET completed=? WHERE id=?", completed, id)
    return err
}

func DeleteTodo(id int) error {
    _, err := db.Exec("DELETE FROM items WHERE id=?", id)
    return err
}

func Close() {
    db.Close()
}

func toNullInt64(num int64) sql.NullInt64 {
    if num == 0 {
        return sql.NullInt64{Int64: 0, Valid: false}
    }
    return sql.NullInt64{Int64: num, Valid: true}
}

func fromNullInt64(num sql.NullInt64) int64 {
    if num.Valid {
        return num.Int64
    }
    return 0
}
