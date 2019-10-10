package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json"
)

type jsonError struct {
	Error string
}

// BadRequest logs a warning level message containing both the error and message.
// A JSON error of the message with a http status bad request is returned.
func BadRequest(err error, message string) (interface{}, int, error) {
	log.Warningln("bad request occurred with message '", err, ": ", message, "'")
	return nil, http.StatusBadRequest, errors.New(message)
}

// UnAuthorized will return an access denied JSON error message with
// http status unauthorized.
func UnAuthorized() (interface{}, int, error) {
	return nil, http.StatusUnauthorized, errors.New("access denied")
}

// NotFound will return a not found JSON error with http status not found
func NotFound() (interface{}, int, error) {
	return nil, http.StatusNotFound, errors.New("not found")
}

// InternalServerError logs a warning level message containing both the error and message.
// A JSON error of the message with a http status internal server error is returned.
func InternalServerError(err error, message string) (interface{}, int, error) {
	log.Error("internal server error occurred with message '", err, ": ", message, "'")
	return nil, http.StatusInternalServerError, errors.New(message)
}

// HandleError will send a JSON error message and set the http status.
func HandleError(w http.ResponseWriter, status int, err error) {
	j := jsonError{Error: err.Error()}
	w.Header().Set(contentType, jsonContentType)
	w.WriteHeader(status)

	e := json.NewEncoder(w).Encode(j)
	if e != nil {
		log.Error("error attempting to encode error message '", err, "'")
	}
}

// Handle attempts to simplify route handlers by using thunks as route handlers, thus removing any direct sending
// of responses from within each route and handling errors here instead.
// Thunks should return: <model to be serialized>, <http status>, <error>
// Models are seralized to JSON
func Handle(method func(r *http.Request) (interface{}, int, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, s, e := method(r)
		if e != nil {
			HandleError(w, s, e)
			return
		}
		if b != nil {
			w.Header().Set(contentType, jsonContentType)
			w.WriteHeader(s)
			err := json.NewEncoder(w).Encode(b)
			if err != nil {
				log.Error("error attempting to encode error message '", err, "'")
			}
		} else {
			w.WriteHeader(s)
		}
	}
}
