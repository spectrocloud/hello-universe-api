package internal

import "database/sql"

type Counter struct {
	Id      int           `json:"id" db:"id"`
	Date    *sql.NullTime `json:"date" db:"date"`
	Browser string        `json:"browser" db:"browser"`
	Os      string        `json:"os" db:"os"`
}
