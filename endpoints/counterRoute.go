// Copyright (c) HashiCorp, Inc.

package endpoints

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mileusna/useragent"
	"github.com/rs/zerolog/log"
	"spectrocloud.com/hello-universe-api/internal"
)

// NewHandlerContext returns a new CounterRoute with a database connection.
func NewCounterHandlerContext(db *sqlx.DB, ctx context.Context, authorization bool) *CounterRoute {
	return &CounterRoute{db, ctx, authorization}
}

func (route *CounterRoute) CounterHTTPHandler(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msg("POST request received. Incrementing counter.")
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	var payload []byte

	if route.authorization && request.Method != "OPTIONS" {
		validation := internal.ValidateToken(request.Header.Get("Authorization"))
		if !validation {
			log.Info().Msg("Invalid token.")
			http.Error(writer, "Invalid credentials.", http.StatusUnauthorized)
			return
		}
	}

	switch request.Method {
	case "POST":
		value, err := route.postHandler(request)
		if err != nil {
			log.Debug().Msg("Error incrementing counter.")
			http.Error(writer, "Error incrementing counter.", http.StatusInternalServerError)
		}
		writer.WriteHeader(http.StatusCreated)
		payload = value
	case "GET":
		value, err := route.getHandler(request)
		if err != nil {
			log.Debug().Msg("Error getting counter value.")
			http.Error(writer, "Error getting counter value.", http.StatusInternalServerError)
		}
		writer.WriteHeader(http.StatusOK)
		payload = value
	case "OPTIONS":
		log.Debug().Msg("OPTIONS request received.")
		log.Debug().Interface("request", request.Header).Msg("Request received.")
		writer.WriteHeader(http.StatusOK)
	default:
		log.Debug().Msg("Invalid request method.")
		http.Error(writer, "Invalid request method.", http.StatusMethodNotAllowed)
	}

	_, err := writer.Write([]byte(payload))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response to the counter endpoint.")
	}
}

// postHandler increments the counter in the database.
func (route *CounterRoute) postHandler(r *http.Request) ([]byte, error) {
	currentTime := time.Now().UTC()
	ua := useragent.Parse(r.UserAgent())
	browser := ua.Name
	os := ua.OS
	transaction, err := route.DB.BeginTx(route.ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error beginning transaction.")
		return []byte{}, err
	}
	sqlQuery := `INSERT INTO counter(date,browser,os) VALUES ($1, $2, $3)`
	_, err = transaction.ExecContext(route.ctx, sqlQuery, currentTime, browser, os)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting counter value.")
		log.Debug().Msgf("SQL query: %s", sqlQuery)
		return []byte{}, err
	}
	log.Info().Msg("Counter incremented in database.")
	getNewCountQuery := `SELECT COUNT(*) AS total FROM counter`
	var databaseTotal sql.NullInt64
	result := transaction.QueryRowContext(route.ctx, getNewCountQuery)
	err = result.Scan(&databaseTotal)
	if err != nil {
		log.Error().Err(err).Msg("Error scanning counter value.")
		return []byte{}, err
	}
	if !databaseTotal.Valid {
		log.Error().Err(err).Msg("Counter value is null.")
		return []byte{}, err
	}
	counterSummary := counterSummary{Total: databaseTotal.Int64}
	err = transaction.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Error committing transaction.")
		err = transaction.Rollback()
		if err != nil {
			log.Error().Err(err).Msg("Error rolling back transaction.")
			panic(err)
		}
		return []byte{}, err
	}
	payload, err := json.MarshalIndent(counterSummary, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling counterSummary struct into JSON.")
		return []byte{}, err
	}
	return payload, nil
}

// getHandler returns the current counter value from the database as a JSON object.
func (route *CounterRoute) getHandler(r *http.Request) ([]byte, error) {
	sqlQuery := `SELECT COUNT(*) AS total FROM counter`
	var counterSummary counterSummary
	err := route.DB.GetContext(route.ctx, &counterSummary, sqlQuery)
	if err != nil {
		log.Error().Err(err).Msg("Error getting counter value.")
		log.Debug().Msgf("SQL query: %s", sqlQuery)
		return []byte{}, err
	}
	log.Info().Msg("Counter value retrieved from database.")
	payload, err := json.MarshalIndent(counterSummary, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling counterSummary struct into JSON.")
		return []byte{}, err
	}
	return payload, nil
}
