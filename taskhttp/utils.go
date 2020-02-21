package taskhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func decode(r *http.Request, v interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read from source during decode: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}

func internalServerError(w http.ResponseWriter) {
	respondJSONError(w, http.StatusInternalServerError, "internal server error")
}

func respondJSONError(w http.ResponseWriter, code int, err string) {
	respondJSONErrors(w, code, []string{err})
}

func respondJSONErrors(w http.ResponseWriter, code int, errs []string) {
	data := map[string]interface{}{
		"code":   code,
		"status": http.StatusText(code),
		"errors": errs,
	}

	respondJSON(w, code, data)
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		panic("json marshal error on response, this is a bug")
	}

	w.WriteHeader(code)
	w.Write(out)
}
