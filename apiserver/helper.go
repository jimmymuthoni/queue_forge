package apiserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
)


type ErrorWithStatus struct {
	status  int
	err 	error
}


func (e *ErrorWithStatus) Error() string {
	return e.err.Error()
}


func handler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			status := http.StatusInternalServerError
			msg := http.StatusText(status)
			if e, ok := err.(*ErrorWithStatus); ok {
				status = e.status
				msg = http.StatusText(e.status)
				if status == http.StatusBadRequest {
					msg = e.err.Error()
				}				

			}
			slog.Error("error executing handler", "error", err, "status", status, "message", msg)
			w.WriteHeader(status)
			if err := json.NewEncoder(w).Encode(APIResponse[struct{}]{
				Message: msg,
			}); err != nil {
				slog.Error("error encoding response", "error", err)
			}
		}
	}
}