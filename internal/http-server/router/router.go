package router

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/magneless/todo-app/internal/http-server/handlers/auth"
	"github.com/magneless/todo-app/internal/http-server/handlers/item"
	"github.com/magneless/todo-app/internal/http-server/handlers/list"
	"github.com/magneless/todo-app/internal/http-server/middleware/authorization"
	mwLogger "github.com/magneless/todo-app/internal/http-server/middleware/logger"
	"github.com/magneless/todo-app/internal/models"
)

type Auth interface {
	CreateUser(name, username, hash_password string) (int64, error)
	GetUser(username, hash_password string) (*models.User, error)
	AddRefreshToken(refresh_token string, expires_at time.Time, username string) error
}

type Repository interface {
	Auth
}

func New(log *slog.Logger, repo Repository) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", auth.SignUp(log, repo))
		r.Post("/sign-in", auth.SignIn(log, repo))
	})

	router.Route("/api", func(r chi.Router) {
		r.Use(authorization.New(log))

		r.Post("/", list.CreateList())
		r.Get("/", list.GetAllLists())
		r.Get("/{id}", list.GetListById())
		r.Put("/{id}", list.UpdateList())
		r.Delete("/{id}", list.DeleteList())

		r.Route("/{id}/items", func(r chi.Router) {
			r.Post("/", item.CreateItem())
			r.Get("/", item.GetAllItems())
			r.Get("/{item_id}", item.GetItemById())
			r.Put("/{item_id}", item.UpdateItem())
			r.Delete("/{item_id}", item.DeleteItem())
		})
	})

	return router
}
