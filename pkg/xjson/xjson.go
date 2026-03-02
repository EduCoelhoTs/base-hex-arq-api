package xjson

import (
	"encoding/json"
	"io"
	"net/http"
)

func Decode[T any](input io.Reader, output *T) error {
	if err := json.NewDecoder(input).Decode(output); err != nil {
		return err
	}

	return nil
}

func ReponseHttp(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

func ResponseHttpError(w http.ResponseWriter, statusCode int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}
