// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: MPL-2.0

package endpoints_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"spectrocloud.com/hello-universe-api/endpoints"
)

func TestNewCounterHandlerContext(t *testing.T) {
	container, err := CreatePostgresTestContainer()
	if err != nil {
		t.Errorf("Error creating PostgresTestContainer: %s", err)
	}

	db, err := StartTestDB(container)
	if err != nil {
		t.Errorf("Expected a new database connection, but got %s", err)
	}

	defer CleanUpTestContainer(container)

	ctx := context.Background()
	authorization := true

	counter := endpoints.NewCounterHandlerContext(db, ctx, authorization)
	if counter == nil {
		t.Errorf("Expected a new CounterRoute, but got nil")
	}

	if counter != nil {

		if counter.Ctx != ctx {
			t.Errorf("Expected context to be %v, but got %v", ctx, counter.Ctx)
		}

		if counter.Authorization != authorization {
			t.Errorf("Expected authorization to be %v, but got %v", authorization, counter.Authorization)
		}

	}

}

func TestCounterHTTPHandlerGETAllPages(t *testing.T) {

	page := "test"
	container, err := CreatePostgresTestContainer()
	if err != nil {
		t.Errorf("Error creating PostgresTestContainer: %s", err)
	}

	db, err := StartTestDB(container)
	if err != nil {
		t.Errorf("Expected a new database connection, but got %s", err)
	}

	defer CleanUpTestContainer(container)

	sqlQuery := `INSERT INTO counter(page,date,browser,os) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlQuery, page, time.Now(), "Chrome", "Windows")
	if err != nil {
		t.Errorf("Error inserting into counter table: %s", err)
	}

	counter := endpoints.NewCounterHandlerContext(db, context.Background(), false)

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

	var result endpoints.CounterSummary

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

func TestCounterHTTPHandlerGETOnePage(t *testing.T) {

	page := "test"
	container, err := CreatePostgresTestContainer()
	if err != nil {
		t.Errorf("Error creating PostgresTestContainer: %s", err)
	}

	db, err := StartTestDB(container)
	if err != nil {
		t.Errorf("Expected a new database connection, but got %s", err)
	}

	defer CleanUpTestContainer(container)

	sqlQuery := `INSERT INTO counter(page,date,browser,os) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlQuery, page, time.Now(), "Chrome", "Windows")
	if err != nil {
		t.Errorf("Error inserting into counter table: %s", err)
	}

	counter := endpoints.NewCounterHandlerContext(db, context.Background(), false)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("v1/counter/%s", page), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("page", page)

	handler := http.HandlerFunc(counter.CounterHTTPHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var result endpoints.CounterSummary

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

	page := "test"
	container, err := CreatePostgresTestContainer()
	if err != nil {
		t.Errorf("Error creating PostgresTestContainer: %s", err)
	}

	db, err := StartTestDB(container)
	if err != nil {
		t.Errorf("Expected a new database connection, but got %s", err)
	}

	defer CleanUpTestContainer(container)

	counter := endpoints.NewCounterHandlerContext(db, context.Background(), false)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("v1/counter/%s", page), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("page", page)

	handler := http.HandlerFunc(counter.CounterHTTPHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var result endpoints.CounterSummary

	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("Error unmarshalling response: %s", err)
	}

	if result.Total < 0 {
		t.Errorf("handler total returned unexpected body: got %v want %s",
			result.Total, "larger than zero")
	}

	sqlQuery := `SELECT COUNT(*) AS total FROM counter WHERE page = $1`
	var counterSummary endpoints.CounterSummary
	err = db.GetContext(context.Background(), &counterSummary, sqlQuery, page)
	if err != nil {
		log.Error().Err(err).Msg("Error getting counter value.")
		log.Debug().Msgf("SQL query: %s", sqlQuery)
	}

	if counterSummary.Total < 1 {
		t.Errorf("handler returned unexpected body: got %v want %s",
			counterSummary.Total, "larger than zero")
	}
}

func TestCounterHTTPHandlerPOSTNoPage(t *testing.T) {

	container, err := CreatePostgresTestContainer()
	if err != nil {
		t.Errorf("Error creating PostgresTestContainer: %s", err)
	}

	db, err := StartTestDB(container)
	if err != nil {
		t.Errorf("Expected a new database connection, but got %s", err)
	}

	defer CleanUpTestContainer(container)

	counter := endpoints.NewCounterHandlerContext(db, context.Background(), false)

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

	var result endpoints.CounterSummary
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("Error unmarshalling response: %s", err)
	}

	if result.Total < 0 {
		t.Errorf("handler total returned unexpected body: got %v want %s",
			result.Total, "larger than zero")
	}

	sqlQuery := `SELECT COUNT(*) AS total FROM counter`
	var counterSummary endpoints.CounterSummary
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
