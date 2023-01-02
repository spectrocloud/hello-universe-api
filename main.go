package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"spectrocloud.com/hello-universe-api/endpoints"
	"spectrocloud.com/hello-universe-api/internal"
)

const (
	db_driver string = "postgres"
)

var (
	dbName           string
	dbUser           string
	dbPassword       string
	dbHost           string
	dbPort           int64
	globalTraceLevel string
	globalDb         *sqlx.DB
	globalHost       string
	globalPort       string
	globalHostURL    string = globalHost + ":" + globalPort
)

func init() {
	globalTraceLevel = internal.Getenv("TRACE", "INFO")
	port := internal.Getenv("PORT", "3000")
	host := internal.Getenv("HOST", "localhost")
	globalHostURL = host + ":" + port

	internal.InitLogger(globalTraceLevel)
	dbName = internal.Getenv("DB_NAME", "counter")
	dbUser = internal.Getenv("DB_USER", "postgres")
	dbHost = internal.Getenv("DB_HOST", "localhost")
	dbPassword = internal.Getenv("DB_PASSWORD", "password")
	dbPort = internal.StringToInt64(internal.Getenv("DB_PORT", "5432"))
	db, err := sqlx.Open(db_driver, fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s connect_timeout=5 sslmode=disable",
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
	))
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	db.SetConnMaxIdleTime(45 * time.Second)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(1 * time.Minute)

	log.Debug().Msg("Checking database connection...")
	err = db.Ping()
	if err != nil {
		log.Debug().Msg("Database is not available")
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	globalDb = db
}

func main() {
	ctx := context.Background()
	counterRoute := endpoints.NewHandlerContext(globalDb, ctx)

	http.HandleFunc(internal.ApiPrefix+"counter", counterRoute.HTTPHandler)

	log.Info().Msg("Server is configured for port 3000")
	log.Info().Msgf("Trace level set to: %s", globalTraceLevel)
	log.Info().Msg("Starting client Application")
	err := http.ListenAndServe(globalHostURL, nil)
	if err != nil {
		log.Debug().Err(err)
		log.Fatal().Msg("There's an error with the server")
	}

}
