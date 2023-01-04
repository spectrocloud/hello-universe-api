package endpoints

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mileusna/useragent"
	"github.com/rs/zerolog/log"
)

// NewHandlerContext returns a new CounterRoute with a database connection.
func NewCounterHandlerContext(db *sqlx.DB, ctx context.Context) *CounterRoute {
	return &CounterRoute{db, ctx}
}

func (route *CounterRoute) CounterHTTPHandler(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msg("POST request received. Incrementing counter.")
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	var payload []byte

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
	sqlQuery := `INSERT INTO counter(date,browser,os) VALUES ($1, $2, $3)`
	_, err := route.DB.ExecContext(route.ctx, sqlQuery, currentTime, browser, os)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting counter value.")
		log.Debug().Msgf("SQL query: %s", sqlQuery)
		return []byte{}, err
	}
	log.Info().Msg("Counter incremented in database.")
	return []byte{}, nil
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
