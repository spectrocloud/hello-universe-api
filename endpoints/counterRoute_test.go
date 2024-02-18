// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: MPL-2.0

package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// startDB returns a new database connection to the counter database.
// A local database is required to run the tests.
func startDB() (*sqlx.DB, error) {
	dbUser := "postgres"
	dbPassword := "password"
	dbName := "counter"
	host := "localhost"
	dbEncryption := "disable"

	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s connect_timeout=5 sslmode=%s",
		host,
		5432,
		dbName,
		dbUser,
		dbPassword,
		dbEncryption,
	))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err

}

func TestNewCounterHandlerContext(t *testing.T) {

	db, err := startDB()
	if err != nil {
		t.Errorf("Expected a new database connection, but got %s", err)
	}

	ctx := context.Background()
	authorization := true

	counter := NewCounterHandlerContext(db, ctx, authorization)
	if counter == nil {
		t.Errorf("Expected a new CounterRoute, but got nil")
	}

	if counter != nil {

		if counter.ctx != ctx {
			t.Errorf("Expected context to be %v, but got %v", ctx, counter.ctx)
		}

		if counter.authorization != authorization {
			t.Errorf("Expected authorization to be %v, but got %v", authorization, counter.authorization)
		}

	}

}

func TestCounterHTTPHandlerGET(t *testing.T) {

	db, err := startDB()
	if err != nil {
		t.Errorf("Expected a new database connection, but got %s", err)
	}

	sqlQuery := `INSERT INTO counter(date,browser,os) VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlQuery, time.Now(), "Chrome", "Windows")
	if err != nil {
		t.Errorf("Error inserting into counter table: %s", err)
	}

	counter := NewCounterHandlerContext(db, context.Background(), false)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "v1/counter", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(counter.CounterHTTPHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var result counterSummary

	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("Error unmarshalling response: %s", err)
	}

fmt.Println(result)
	

	if result.Total == 0 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			result.Total, 0)
	}
}


func TestCounterHTTPHandlerPOST(t *testing.T) {
	
	db, err := startDB()
	if err != nil {
		t.Errorf("Expected a new database connection, but got %s", err)
	}

	counter := NewCounterHandlerContext(db, context.Background(), false)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "v1/counter", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(counter.CounterHTTPHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var result counterSummary

	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("Error unmarshalling response: %s", err)
	}
	


	if result.Total < 0 {
		t.Errorf("handler returned unexpected body: got %v want %s",
			result.Total, "larger than zero")
	}

	if len(result.Counts) < 0 {
		t.Errorf("handler returned unexpected body: got %v want %s",
			len(result.Counts), "larger than zero")
	}

	sqlQuery := `SELECT COUNT(*) AS total FROM counter`
	var counterSummary counterSummary
	err = db.GetContext(context.Background(), &counterSummary, sqlQuery)
	if err != nil {
		log.Error().Err(err).Msg("Error getting counter value.")
		log.Debug().Msgf("SQL query: %s", sqlQuery)
	}

	if counterSummary.Total < 1 {
		t.Errorf("handler returned unexpected body: got %v want %s",
			counterSummary.Total, "larger than zero")
	}
}
