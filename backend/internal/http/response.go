package http

import (
	"encoding/json"
	"net/http"
)

type responseEnvelope struct {
	Data  any            `json:"data"`
	Error *errorEnvelope `json:"error"`
}

type errorEnvelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(responseEnvelope{Data: data, Error: nil})
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(responseEnvelope{
		Data: nil,
		Error: &errorEnvelope{
			Code:    code,
			Message: message,
		},
	})
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
