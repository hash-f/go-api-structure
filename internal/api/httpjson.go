package api

import (
	"encoding/json"
	"errors"
	"fmt" // For methodNotAllowedResponse
	"io"  // For io.EOF and io.ErrUnexpectedEOF
	"log/slog"
	"net/http"
	"strings" // For checking unknown field errors
)

// encode writes a JSON response with the given status code and data.
// If status is http.StatusNoContent, it writes a response with no body.
// Otherwise, it attempts to JSON encode data. If data is a nil pointer/interface,
// it will be encoded as JSON 'null'.
// It logs an error if JSON marshaling fails.
func encode[T any](w http.ResponseWriter, r *http.Request, status int, data T) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}

	// If data is an interface type and is nil, json.NewEncoder will write "null".
	// If data is a pointer type and is nil, json.NewEncoder will write "null".
	// This is standard JSON behavior.
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// If encoding fails, log the error. The response status code is already sent.
		// It's hard to send a different error to the client at this point if headers are written.
		logger, _ := r.Context().Value(GetLoggerKey()).(*slog.Logger) // Use GetLoggerKey, handle potential nil logger
		if logger != nil {
			logger.Error("failed to encode JSON response", "error", err)
		} else {
			// Fallback logger if not found in context (should not happen with middleware)
			slog.Default().Error("failed to encode JSON response and no logger in context", "error", err)
		}
		// Attempt to send a generic error if possible, though likely too late.
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Define a context key type for the logger to avoid collisions.
// This will be used by middleware to inject the logger and by handlers to retrieve it.
type contextKey string

const loggerKey = contextKey("logger")

// GetLoggerKey returns the key used to store the logger in the request context.
// This is exported so that middleware (e.g. in server package) can use it to set the logger.
func GetLoggerKey() any {
	return loggerKey
}

const maxRequestSize = 1_048_576 // 1 MB

// decode reads a JSON request body into the provided destination struct `dst`.
// It enforces a maximum request body size (maxRequestSize).
// It handles various decoding errors, returning specific error types or messages
// that can be translated into appropriate HTTP error responses by the caller.
func decode[T any](w http.ResponseWriter, r *http.Request, dst T) error {
	// Set a maximum body size to prevent abuse.
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	// Use a DisallowUnknownFields decoder to prevent unexpected fields.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Check for unknown fields error. This error message is specific to json.Decoder.DisallowUnknownFields().
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			// Remove surrounding quotes from fieldName if present
			fieldName = strings.Trim(fieldName, `"`)
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Check for body too large error. This error message is specific to http.MaxBytesReader.
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxRequestSize)

		case errors.As(err, &invalidUnmarshalError):
			// This error occurs if a nil pointer is passed to Decode, which is a panic.
			// It's more of a server-side programming error.
			panic(err) // Or return a generic server error

		default:
			return err
		}
	}

	// Call Decode again, using a pointer to an empty anonymous struct as the destination.
	// If the request body only contained a single JSON value this will return an io.EOF error.
	// If there is anything else, the error will be different and we know there is more than one JSON value.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// envelope is a helper type for wrapping JSON responses, allowing for consistent structure.
// For example, data can be wrapped as {"data": ...} or errors as {"error": ...}.
type envelope map[string]any

// errorResponse sends a JSON formatted error message with the given status code.
func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	envelope := envelope{"error": message}
	encode(w, r, status, envelope)
}

// serverErrorResponse sends a 500 Internal Server Error response.
// It logs the error before sending the generic message to the client.
func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger, _ := r.Context().Value(GetLoggerKey()).(*slog.Logger)
	if logger != nil {
		logger.Error("internal server error", "error", err.Error())
	} else {
		slog.Default().Error("internal server error and no logger in context", "error", err.Error())
	}

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

// badRequestResponse sends a 400 Bad Request response.
func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// notFoundResponse sends a 404 Not Found response.
func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse sends a 405 Method Not Allowed response.
func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// failedValidationResponse sends a 422 Unprocessable Entity response.
func failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// TODO: Add more specific error responses as needed (e.g., authentication errors).

// decodeAndValidate decodes the JSON request body into dst and then validates dst.
// dst must be a pointer to a struct that implements the Validator interface.
// If decoding fails, it sends an appropriate error response (e.g., 400 Bad Request).
// If validation fails, it sends a 422 Unprocessable Entity response with validation errors.
// Returns true if both decoding and validation are successful, false otherwise.
func decodeAndValidate[T Validator](w http.ResponseWriter, r *http.Request, dst T) bool {
	err := decode(w, r, dst) // decode expects a pointer, T should be a pointer type that implements Validator
	if err != nil {
		// decode function already returns specific errors that can be mapped to badRequestResponse
		// or serverErrorResponse if it's a more generic error.
		// For simplicity here, we'll assume most decode errors are client-side.
		// A more robust error handling might inspect err further.
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError), errors.Is(err, io.ErrUnexpectedEOF), errors.As(err, &unmarshalTypeError), errors.Is(err, io.EOF), strings.HasPrefix(err.Error(), "json: unknown field"), err.Error() == "http: request body too large", err.Error() == "body must only contain a single JSON value":
			badRequestResponse(w, r, err)
		default:
			serverErrorResponse(w, r, err) // For other unexpected errors during decode
		}
		return false
	}

	// Perform validation
	if validationErrors := dst.Valid(); validationErrors != nil {
		failedValidationResponse(w, r, validationErrors)
		return false
	}

	return true
}
