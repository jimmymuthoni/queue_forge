package apiserver

import (
	"encoding/json"
	"fmt"
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

func NewErrWithStatus(status int, err error) *ErrorWithStatus {
	return &ErrorWithStatus{
		status: status,
		err: err,
	}
}


//error response helper
func handler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			status := http.StatusInternalServerError
			msg := http.StatusText(status)
			if e, ok := err.(*ErrorWithStatus); ok {
				status = e.status
				msg = http.StatusText(e.status)
				if status == http.StatusBadRequest || status == http.StatusConflict {
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


//encode helper
func encode[T any](v T, status int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encoding reposnse: %w", err)
	}
	return nil
}


type Varidator interface {
	Validate() error
}
//decode helper
func decode[T Varidator](r *http.Request) (T,error) {
	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return t, fmt.Errorf("decoding request body: %w", err)
	}
	if err := t.Validate(); err != nil {
		return t, err
	}
	return t, nil
}