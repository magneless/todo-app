package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/magneless/todo-app/internal/http-server/handlers/auth"
	"github.com/magneless/todo-app/internal/http-server/handlers/item"
	"github.com/magneless/todo-app/internal/http-server/handlers/list"
)


func InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", auth.SignUp())
		r.Post("/sign-in", auth.SignIn())
	})

	router.Route("/api", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
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
	})

	return router
}
