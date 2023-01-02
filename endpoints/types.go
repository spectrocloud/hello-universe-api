package endpoints

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Counter struct {
	Id      int           `json:"id" db:"id"`
	Date    *sql.NullTime `json:"date" db:"date"`
	Browser string        `json:"browser" db:"browser"`
	Os      string        `json:"os" db:"os"`
}

type CounterRoute struct {
	DB  *sqlx.DB
	ctx context.Context
}

type counterSummary struct {
	Total  int       `json:"total" db:"total"`
	Counts []Counter `json:"counts" db:"counts"`
}

type HeatlhRoute struct {
	ctx context.Context
}
