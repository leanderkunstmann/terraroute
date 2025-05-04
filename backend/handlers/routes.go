package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const basePathV1 = "/api/v1"

type Handler interface {
	Register(r *mux.Router)
}

// newErrorResponse is a helper function to create a JSON error response.
// It sets the Content-Type header to application/json and writes the error message
// and status code to the response writer.
// If the error implements the Code() method, it uses that to determine the status code.
func newErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	code := codeFromError(err, statusCode)
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode error response: %v", err), http.StatusInternalServerError)
	}
}

// codeFromError is a helper function to extract the status code from an error.
// It checks if the error implements the Code() method, and if so, returns that code.
// If the error is nil or does not implement the Code() method, it returns the default code.
func codeFromError(err error, defaultCode int) int {
	if err == nil {
		return defaultCode
	}

	if customErr, ok := err.(interface{ Code() int }); ok {
		return customErr.Code()
	}

	if customErr, ok := err.(interface{ Unwrap() error }); ok {
		return codeFromError(customErr.Unwrap(), defaultCode)
	}

	if customErr, ok := err.(interface{ Unwrap() []error }); ok {
		errs := customErr.Unwrap()
		for _, e := range errs {
			if code := codeFromError(e, defaultCode); code != defaultCode {
				// Return the first non-internal server error code
				return code
			}
		}
	}

	return defaultCode
}
