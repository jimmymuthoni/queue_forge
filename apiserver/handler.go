package apiserver

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
	"uuid"
)


type SignUpRequest struct {
	Email    string  `json:"email"`
	Password string	 `json:"password"`
}


type TokenRefreshResponse struct {
	AccessToken 	string `json:"access_token"`
	RefreshToken	string `json:"refresh_token"`
}


//response for signinHandler
type SigninResponse struct {
	AccessToken 		string 	`json:"access_token"`
	RefreshToken 	string  `json:"refresh_token"`
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


//payload that hold user info
type SignInRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
}

//validating user info
func (r SignInRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

//signinHadler sign in user and issue token
func (s *ApiServer) signinHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		req, err := decode[SignInRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, err)
		}

		//check if user exist and retrieve
		user, err := s.store.Users.FindUserByEmail(r.Context(), req.Email)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		//verify is password matches
		if err := user.ComparePassword(req.Password); err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}

		//giving token to verified user
		tokenPair, err:= s.jwtManager.GenerateTokenPair(user.Id)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		_, err = s.store.RefreshTokenStore.DeleteUserTokens(r.Context(), user.Id)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		_, err = s.store.RefreshTokenStore.Create(r.Context(), user.Id, tokenPair.RefreshToken)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if err := encode(APIResponse[SigninResponse]{
			Data: &SigninResponse{
				AccessToken: tokenPair.AccessToken.Raw,
				RefreshToken: tokenPair.RefreshToken.Raw,
			},
		},http.StatusOK, w); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		
		return nil
	})
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

//validating the refresh token
func (r TokenRefreshRequest) Validate() error {
	if r.RefreshToken == " " {
		return errors.New("refresh token iis required")
	}
	return nil
}


func tokenRefreshHandler(s *ApiServer) http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		req, err := decode[TokenRefreshRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, err)
		}

		currentRefreshToken, err := s.jwtManager.Parse(req.RefreshToken)
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}

		userIdstr, err := currentRefreshToken.Claims.GetSubject()
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}

		userId, err := uuid.Parse(userIdstr)
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}

		currentRefreshTokenRecord, err := s.store.RefreshTokenStore.ByPrimaryKey(r.Context(), uuid.UUID(userId), currentRefreshToken)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, sql.ErrNoRows){
				status = http.StatusUnauthorized
			}
			return NewErrWithStatus(status, err)
		}

		if currentRefreshTokenRecord.ExpiresAt.Before(time.Now()){
			return NewErrWithStatus(http.StatusUnauthorized, fmt.Errorf("refrsh token expired"))
		}

		tokenPair, err := s.jwtManager.GenerateTokenPair(userId)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)

		}

		if _, err := s.store.RefreshTokenStore.DeleteUserTokens(r.Context(), uuid.UUID(userId)); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)

		}

		if _, err := s.store.RefreshTokenStore.Create(r.Context(), userId, tokenPair.RefreshToken); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		if err := encode(APIResponse[TokenRefreshResponse]{
			Data: &TokenRefreshResponse{
				AccessToken: tokenPair.AccessToken.Raw,
				RefreshToken: tokenPair.AccessToken.Raw,
			},
		}, http.StatusOK, w); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		
		return nil

	})
}