package endpoints

import (
	"net/http"
)

type Route interface {
	HTTPHandler(writer http.ResponseWriter, request *http.Request)
}
