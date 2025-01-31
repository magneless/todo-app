package auth

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/magneless/todo-app/internal/lib/api/response"
	"github.com/magneless/todo-app/internal/lib/hashing"
	jwt_token "github.com/magneless/todo-app/internal/lib/jwt"
	"github.com/magneless/todo-app/internal/lib/logger/sl"
	"github.com/magneless/todo-app/internal/models"
	"github.com/magneless/todo-app/internal/storage"
)

type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpRequest struct {
	SignInRequest
	Name string `json:"name" validate:"required"`
}

type UserCreater interface {
	CreateUser(name, username, hash_password string) (int64, error)
}

type UserGetter interface {
	GetUser(username, hash_password string) (*models.User, error)
}

func SignUp(log *slog.Logger, userCreater UserCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.signup"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SignUpRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.InternalError())

			return
		}

		log.Info("request body decoded ", slog.Any("request", map[string]string{
			"Name":     req.Name,
			"Username": req.Username,
			"Password": "******",
		}))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		passwordHash, err := hashing.HashPassword(req.Password)
		if err != nil {
			log.Error("failed to hash password", sl.Err(err))

			render.JSON(w, r, resp.Error("not allowed password"))

			return
		}

		id, err := userCreater.CreateUser(req.Name, req.Username, passwordHash)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("username", req.Username))

			render.JSON(w, r, resp.Error("user already exists"))

			return
		}
		if err != nil {
			log.Error("failed to sign up user", sl.Err(err))

			render.JSON(w, r, resp.InternalError())

			return
		}

		log.Info("user sign up", slog.Int64("id", id))

		render.JSON(w, r, resp.OK())
	}
}

func SignIn(log *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.signin"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SignInRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.InternalError())

			return
		}

		log.Info("request body decoded ", slog.Any("request", map[string]string{
			"Username": req.Username,
			"Password": "******",
		}))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		passwordHash, err := hashing.HashPassword(req.Password)
		if err != nil {
			log.Error("failed to hash password", sl.Err(err))

			render.JSON(w, r, resp.Error("not allowed password"))

			return
		}

		user, err := userGetter.GetUser(req.Username, passwordHash)
		if err != nil {
			log.Error("failed to sign in user", sl.Err(err))

			render.JSON(w, r, resp.Error("wrong username or password"))

			return
		}

		accessToken, err := jwt_token.GenerateAccessToken(user.Username)
		if err != nil {
			log.Error("failed to generate access token", sl.Err(err))

			render.JSON(w, r, resp.InternalError())

			return
		}

		refreshToken, err := jwt_token.GenerateRefreshToken(user.Username)
		if err != nil {
			log.Error("failed to generate refresh token", sl.Err(err))

			render.JSON(w, r, resp.InternalError())

			return
		}

		log.Info("user sign in", slog.String("username", req.Username))

		render.JSON(w, r, resp.OKWithData(map[string]string{
			"Access Token":  accessToken,
			"Refresh Token": refreshToken,
		}))
	}
}
