package main

import (
	"log"
	"net/http"
	"os"

	"github.com/magneless/todo-app/internal/config"
	"github.com/magneless/todo-app/internal/http-server/router"
	reposiory "github.com/magneless/todo-app/internal/repository"
	"github.com/magneless/todo-app/internal/storage/postgre"
)

func main() {
	cfg := config.MustLoad()

	storage, err := postgre.New(cfg.Storage)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	reposiory.New(storage)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router.New(),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	
	if err := srv.ListenAndServe(); err != nil {

	}
}
