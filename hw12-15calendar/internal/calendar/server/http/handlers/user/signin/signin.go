package signin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
)

type UserAuth interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type RequestSignIn struct {
	UserName string `json:"userName" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type ResponseSignIn struct {
	AccessToken string `json:"accessToken"`
	Status      string `json:"status"`
}

// SignIn obtaining an access token.
func New(auth UserAuth) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request RequestSignIn

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("SignIn: read body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("SignIn: unmarshal body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		slog.Debug(
			"Sign in request",
			slog.String("username", request.UserName),
			slog.String("password", request.Password),
		)

		if err := validator.New().Struct(request); err != nil {
			errors := err.(validator.ValidationErrors)
			err := fmt.Errorf("SignIn: invalid request: %s", errors)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		accessToken, err := auth.Authenticate(
			req.Context(),
			request.UserName,
			request.Password,
		)

		// TODO: http.StatusUnauthorized and http.StatusInternalServerError
		if err != nil {
			err := fmt.Errorf("SignIn: authenticate user: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusUnauthorized)
			return
		}

		// Write json response
		response := ResponseSignIn{
			AccessToken: accessToken,
			Status:      "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("SignIn: marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("SignIn: write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Login and token generation were successful")
	}
}
