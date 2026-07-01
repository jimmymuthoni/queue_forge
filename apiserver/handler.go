package apiserver

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	Data    *T		`json:"data"`
	Message string `json:"message,omitempty"`
}

//handler to handle signup of a user
func (s *ApiServer) signupHandler(w http.ResponseWriter, r *http.Request){
	var reg SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := reg.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//getting existing user by email
	existingUser, err := s.store.Users.FindUserByEmail(r.Context(), reg.Email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	if existingUser != nil {
    http.Error(w, "user already exists", http.StatusBadRequest)
    return
	}

	//creating user
	_, err = s.store.Users.CreateUser(r.Context(), reg.Email, reg.Password)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode (APIResponse[struct{}]{
		Message: "Successfullly signup user",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}