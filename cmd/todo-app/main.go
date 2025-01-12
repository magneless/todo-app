package main

import (
	"net/http"

	"github.com/magneless/todo-app/internal/config"
	"github.com/magneless/todo-app/internal/http-server/router"
)

func main() {
	cfg := config.MustLoad()

	router := router.InitRoutes()

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	
	if err := srv.ListenAndServe(); err != nil {

	}
}
