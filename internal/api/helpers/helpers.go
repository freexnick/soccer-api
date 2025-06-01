package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func WriteJSON(w http.ResponseWriter, status uint16, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	w.Write(js)

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1 << 20
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

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
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func ErrorResponse(w http.ResponseWriter, status uint16, message any) {
	env := map[string]any{"error": message}
	if err := WriteJSON(w, status, env, nil); err != nil {
		fmt.Printf("CRITICAL (helpers.ErrorResponse): Could not write error JSON to response: %v (original status: %d, message: %v)\n", err, status, message)
	}
}

func BadRequestResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusBadRequest, err.Error())
}

func ServerErrorResponse(w http.ResponseWriter, internalErr error) {
	fmt.Printf("INTERNAL SERVER ERROR (to be logged with observer): %v\n", internalErr)
	ErrorResponse(w, http.StatusInternalServerError, "The server encountered an unexpected problem.")
}

func UnauthorizedResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Authentication required."
	}
	ErrorResponse(w, http.StatusUnauthorized, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("The requested resource '%s' was not found.", r.URL.Path))
}
