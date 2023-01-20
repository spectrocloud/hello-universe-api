package internal

import (
	"strings"

	"github.com/rs/zerolog/log"
)

// ValidateToken validates the token received from the request is valid. Currently, the only valid token is the AnnoymousToken.
func ValidateToken(token string) bool {

	if token == "" {
		log.Debug().Msg("No token provided.")
		return false
	}

	splitToken := strings.Split(token, "Bearer ")
	requestToken := splitToken[1]

	if requestToken != AnnoymousToken {
		log.Debug().Msg("Invalid token. Received: " + token + "")
		return false
	}

	return true

}
