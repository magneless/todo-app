package authorization

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	resp "github.com/magneless/todo-app/internal/lib/api/response"
	"github.com/magneless/todo-app/internal/lib/jwt"
	"github.com/magneless/todo-app/internal/lib/logger/sl"
)

type ctxKey string

const UsernameKey ctxKey = "username"

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	log = log.With(slog.String("component", "middleware.authorization"))

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("authorization header is missed")

				render.JSON(w, r, resp.Error("authorization header is missed"))

				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Error("invalid authorization header format")

				render.JSON(w, r, resp.Error("invalid authorization header format"))

				return
			}

			tokenString := parts[1]

			username, err := jwt_token.ValidateAccessToken(tokenString)
			if err != nil {
				log.Error("error in token validation", sl.Err(err))

				render.JSON(w, r, resp.Error("invalid or expired access token"))

				return
			}

			ctx := context.WithValue(r.Context(), UsernameKey, username)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
