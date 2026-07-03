package apiserver

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)


type SignUpRequest struct {
	Email    string  `json:"email"`
	Password string	 `json:"password"`
}

func (r SignUpRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil

}

//generic APIResponse can contain data to return from API
type APIResponse[T any] struct{
	Data    *T		`json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

// handler to handle signup of a user
func (s *ApiServer) signupHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {

		// decoding the incoming request
		req, err := decode[SignUpRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, err)
		}

		// getting existing user by email
		existingUser, err := s.store.Users.FindUserByEmail(r.Context(), req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return NewErrWithStatus(http.StatusBadRequest, err)
		}

		if existingUser != nil {
			return NewErrWithStatus(http.StatusConflict, fmt.Errorf("email already registred"))
		}

		// creating user
		_, err = s.store.Users.CreateUser(r.Context(), req.Email, req.Password)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		// encoding API response
		if err := encode(
			APIResponse[struct{}]{
				Message: "successfully signed up user",
			},
			http.StatusCreated,
			w,
		); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		return nil
	})
}