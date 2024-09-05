// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: MPL-2.0

package endpoints

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Counter struct {
	Id      int           `json:"id" db:"id"`
	Page    string        `json:"page" db:"page"`
	Date    *sql.NullTime `json:"date" db:"date"`
	Browser string        `json:"browser" db:"browser"`
	Os      string        `json:"os" db:"os"`
}

type CounterRoute struct {
	DB            *sqlx.DB
	Ctx           context.Context
	Authorization bool
}

type CounterSummary struct {
	Total  int64     `json:"total" db:"total"`
	Counts []Counter `json:"counts,omitempty" db:"counts"`
}

type HealthRoute struct {
	Ctx           context.Context
	Authorization bool
}
