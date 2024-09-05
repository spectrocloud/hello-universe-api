package endpoints_test

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"spectrocloud.com/hello-universe-api/internal"
)

const (
	dbName        string = "counter"
	dbUser        string = "postgres"
	dbPassword    string = "password"
	image         string = "ghcr.io/spectrocloud/hello-universe-db"
	image_version string = "1.1.0"
)

func CreatePostgresTestContainer() (*postgres.PostgresContainer, error) {
	ctx := context.Background()
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
		return nil, err
	}

	return postgresContainer, nil
}

// StartTestDB returns a new database connection to the counter database.
// A local database is required to run the tests.
func StartTestDB(container *postgres.PostgresContainer) (*sqlx.DB, error) {
	ctx := context.Background()
	connection, err := container.ConnectionString(ctx,
		fmt.Sprintf("user=%s", dbUser),
		fmt.Sprintf("password=%s", dbPassword),
		fmt.Sprintf("dbname=%s", dbName),
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

	err = internal.InitDB(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("Expected database initialization, but got %s", err)
	}

	return db, err
}

func CleanUpTestContainer(container *postgres.PostgresContainer) {
	ctx := context.Background()
	if err := container.Terminate(ctx); err != nil {
		log.Fatal().Msgf("failed to terminate container: %s", err)
	}
}
