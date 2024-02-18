package endpoints

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
