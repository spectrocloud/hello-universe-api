package internal

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInitDB(t *testing.T) {

	ctx := context.Background()

	dbName := "counter"
	dbUser := "postgres"
	dbPassword := "password"

	const (
		image          string = "ghcr.io/spectrocloud/hello-universe-db"
		image_version string = "1.0.2"
	)

	postgresContainer, err := postgres.Run(ctx, fmt.Sprintf("%s:%s", image, image_version),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	db, err := startDB(ctx, postgresContainer)
	if err != nil {
		log.Fatalf("failed to start database: %s", err)
	}

	err = InitDB(ctx, db)
	if err != nil {
		t.Errorf("Expected database initailization, but got %s", err)
	}

	// Check if the table was created
	query := `SELECT EXISTS (
		SELECT 1 FROM pg_tables
		WHERE schemaname = 'public'
		AND tablename = 'counter'
	);
	`

	var tableExists bool
	err = db.QueryRowContext(ctx, query).Scan(&tableExists)
	if err != nil {
			t.Errorf("Unable to query the database: %s", err)
	}
	
	
	if !tableExists {
			t.Errorf("Expected table 'counter' to exist, but it does not.")
	}


	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

}

func startDB(ctx context.Context, container *postgres.PostgresContainer) (*sqlx.DB, error) {

	connection, err := container.ConnectionString(ctx, 
		"user=posgres", 
		"password=password", 
		"dbname=counter",
		"sslmode=disable",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get container connection string: %s", err)
	}

	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err

}
