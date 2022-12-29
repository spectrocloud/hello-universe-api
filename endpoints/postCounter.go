package endpoints

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func CounterHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	log.Debug().Msg("POST request received. Incrementing counter.")
	var payload []byte

	if request.Method == "POST" {
		value, err := postHandler(request)
		if err != nil {
			log.Debug().Msg("Error incrementing counter.")
			http.Error(writer, "Error incrementing counter.", http.StatusInternalServerError)
		}
		payload = value
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(payload))
}

func postHandler(r *http.Request) ([]byte, error) {
	// ...
	return []byte{}, nil
}
