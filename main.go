package main

import (
	"net/http"

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
	// 	dbName = "postgres"
	// 	dbUser = "postgres"
	// 	dbHost = getenv("DB_HOST", "localhost")
	// 	dbPassword = "mysecretpassword"
	// 	dbPort = 5432
	// 	applicationHost := getenv("APPLICATION_HOST", "localhost")
	// 	applicationPort := getenv("APPLICATION_PORT", "8080")
	// 	host = fmt.Sprintf("%s:%s", applicationHost, applicationPort)
	// 	db, err := sqlx.Open(db_driver, fmt.Sprintf(
	// 		"host=%s port=%d dbname=%s user=%s password=%s connect_timeout=5 sslmode=disable",
	// 		dbHost,
	// 		dbPort,
	// 		dbName,
	// 		dbUser,
	// 		dbPassword,
	// 	))
	// 	if err != nil {
	// 		log.Printf("Error connecting to database: %v", err)
	// 		log.Fatal(err)
	// 	}

	// 	db.SetConnMaxIdleTime(45 * time.Second)
	// 	db.SetMaxIdleConns(3)
	// 	db.SetConnMaxLifetime(1 * time.Minute)

	// log.Println("Checking database connection...")
	// err = db.Ping()
	//
	//	if err != nil {
	//		log.Print("Database is not available")
	//		log.Print(err)
	//		dbEnabled = false
	//	}
	//
	//	if dbEnabled {
	//		// Set the max value for the random function
	//		err = db.Get(&maxCounter, "SELECT COUNT(*) FROM jokes")
	//		if err != nil {
	//			log.Printf("Unable to get max value: %v", err)
	//		}
	//		globalDb = db
	//	}
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/jokes" {
// 		http.Redirect(w, r, "/jokes", http.StatusFound)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")

// 	switch r.Method {
// 	case "GET":
// 		if dbEnabled {
// 			joke := joke{}
// 			err := globalDb.Get(&joke, "SELECT * FROM jokes WHERE id = $1", randomInt(1, int(maxCounter)))
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			payload := response{
// 				Body: joke.Joke,
// 			}

// 			response, err := json.MarshalIndent(payload, " ", " ")
// 			if err != nil {
// 				log.Println(err.Error())
// 			}
// 			w.Write(response)
// 		} else {
// 			w.Write([]byte("Database is not available"))
// 		}
// 	case "POST":
// 		if dbEnabled {
// 			body, err := io.ReadAll(r.Body)
// 			if err != nil {
// 				log.Printf("Error reading body: %v", err)
// 				http.Error(w, "can't read body", http.StatusBadRequest)
// 				return
// 			}

// 			joke := joke{}
// 			err = json.Unmarshal(body, &joke)
// 			if err != nil {
// 				log.Printf("Error unmarshalling body: %v", err)
// 				http.Error(w, "can't unmarshal body", http.StatusBadRequest)
// 				return
// 			}

// 			_, err = globalDb.Exec("INSERT INTO jokes (joke) VALUES ($1)", joke.Joke)
// 			if err != nil {
// 				log.Printf("Error inserting joke: %v", err)
// 				http.Error(w, "can't insert joke", http.StatusBadRequest)
// 				return
// 			}
// 			maxCounter = maxCounter + 1
// 		} else {
// 			w.Write([]byte("Database is not available"))
// 		}

// 	default:
// 		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
// 	}
// }

func main() {
	http.HandleFunc(internal.ApiPrefix+"counter", endpoints.CounterHandler)
	log.Info().Msg("Server is configured for port 3000")
	log.Info().Msgf("Trace level set to: %s", globalTraceLevel)
	log.Info().Msg("Starting client Application")
	err := http.ListenAndServe(globalHostURL, nil)
	if err != nil {
		log.Debug().Err(err)
		log.Fatal().Msg("There's an error with the server")
	}

}
