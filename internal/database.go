// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
)

// InitDB initializes the database by creating all the required tables.
func InitDB(ctx context.Context, db *sqlx.DB) error {

	sqlStatement := `
	CREATE TABLE IF NOT EXISTS counter (
		id SERIAL PRIMARY KEY,
		date timestamp NOT NULL,
		browser varchar(255),
		os varchar(255)
	);
`
	_, err := db.ExecContext(ctx, sqlStatement)
	if err != nil {
		log.Info().Err(err).Msg("Error initializing the database")
		log.Debug().Msgf("SQL statement: %s", sqlStatement)
		return err
	}
	log.Info().Msg("Database initialized successfully")
	return nil
}
