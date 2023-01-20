package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	dbName              string
	dbUser              string
	dbPassword          string
	dbHost              string
	dbPort              int64
	globalTraceLevel    string
	globalDb            *sqlx.DB
	globalHost          string
	globalPort          string
	globalHostURL       string = globalHost + ":" + globalPort
	globalAuthorization bool
)

func init() {
	globalTraceLevel = strings.ToUpper(internal.Getenv("TRACE", "INFO"))
	internal.InitLogger(globalTraceLevel)
	authorizationEnv := strings.ToUpper(internal.Getenv("AUTHORIZATION", "false"))
	authorization, err := strconv.ParseBool(authorizationEnv)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing authorization")
	}
	globalAuthorization = authorization
	initDB := strings.ToLower(internal.Getenv("DB_INIT", "false"))
	port := internal.Getenv("PORT", "3000")
	host := internal.Getenv("HOST", "0.0.0.0")
	globalHost = host
	globalPort = port
	globalHostURL = host + ":" + port
	dbName = internal.Getenv("DB_NAME", "counter")
	dbUser = internal.Getenv("DB_USER", "postgres")
	dbHost = internal.Getenv("DB_HOST", "0.0.0.0")
	dbEncryption := internal.Getenv("DB_ENCRYPTION", "disable")
	dbPassword = internal.Getenv("DB_PASSWORD", "password")
	dbPort = internal.StringToInt64(internal.Getenv("DB_PORT", "5432"))
	db, err := sqlx.Open(db_driver, fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s connect_timeout=5 sslmode=%s",
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
		dbEncryption,
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

	if initDB == "true" {
		log.Debug().Msg("Initializing database")
		err = internal.InitDB(context.Background(), db)
		if err != nil {
			log.Fatal().Err(err).Msg("Error initializing database")
		}
	}

	globalDb = db
}

func main() {
	ctx := context.Background()
	counterRoute := endpoints.NewCounterHandlerContext(globalDb, ctx, globalAuthorization)
	healthRoute := endpoints.NewHealthHandlerContext(ctx, globalAuthorization)

	http.HandleFunc(internal.ApiPrefix+"counter", counterRoute.CounterHTTPHandler)
	http.HandleFunc(internal.ApiPrefix+"health", healthRoute.HealthHTTPHandler)

	log.Info().Msgf("Server is configured for port %s and listing on %s", globalPort, globalHostURL)
	log.Info().Msgf("Database is configured for %s:%d", dbHost, dbPort)
	log.Info().Msgf("Trace level set to: %s", globalTraceLevel)
	log.Info().Msg("Starting client Application")
	log.Info().Msgf("Authorization is set to: %v", globalAuthorization)
	err := http.ListenAndServe(globalHostURL, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("There's an error with the server")
	}

}
