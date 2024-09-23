package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

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
