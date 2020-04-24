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

type Response struct {
	body   interface{}
	status int
	error  error
}

// BadRequest logs a warning level message containing both the error and message.
// A JSON error of the message with a http status bad request is returned.
func BadRequest(err error, message string) Response {
	log.Warningln("bad request occurred with message '", err, ": ", message, "'")
	return Response{
		body:   nil,
		status: http.StatusBadRequest,
		error:  errors.New(message),
	}
}

// UnAuthorized will return an access denied JSON error message with
// http status unauthorized.
func UnAuthorized() Response {
	return Response{
		body:   nil,
		status: http.StatusUnauthorized,
		error:  errors.New("access denied"),
	}
}

// NotFound will return a not found JSON error with http status not found
func NotFound() Response {
	return Response{
		body:   nil,
		status: http.StatusNotFound,
		error:  errors.New("not found"),
	}
}

// NotContent will return an empty body
func NoContent() Response {
	return Response{
		body:   nil,
		status: http.StatusOK,
		error:  nil,
	}
}

// Ok will return the provided body
func Ok(body interface{}) Response {
	return Response{
		body:   body,
		status: http.StatusOK,
		error:  nil,
	}
}

// InternalServerError logs a warning level message containing both the error and message.
// A JSON error of the message with a http status internal server error is returned.
func InternalServerError(err error, message string) Response {
	log.Error("internal server error occurred with message '", err, ": ", message, "'")
	return Response{
		body:   nil,
		status: http.StatusInternalServerError,
		error:  errors.New(message),
	}
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
func Handle(method func(r *http.Request) Response) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := method(r)
		if response.error != nil {
			HandleError(w, response.status, response.error)
			return
		}
		if response.body != nil {
			w.Header().Set(contentType, jsonContentType)
			w.WriteHeader(response.status)
			err := json.NewEncoder(w).Encode(response.body)
			if err != nil {
				log.Error("error attempting to encode error message '", err, "'")
			}
		} else {
			w.WriteHeader(response.status)
		}
	}
}
