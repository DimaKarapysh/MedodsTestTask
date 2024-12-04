package main

import (
	"MedodsTestTask/app"
	"MedodsTestTask/delivery"
	"MedodsTestTask/iteractors"
	"MedodsTestTask/repository"
	"MedodsTestTask/tools/config"
	"context"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := Run()
	if err != nil {
		log.Fatal(err)

	}
}

func Run() error {
	// Env
	err := config.InitEnv()
	if err != nil {
		return errors.Wrap(err, "Config")
	}

	// Log
	logs, err := app.InitLogs()
	if err != nil {
		return errors.Wrap(err, "Log")
	}

	// Db
	db, err := app.InitDatabase(logs)
	if err != nil {
		return errors.Wrap(err, "DB")
	}

	// Migration
	err = app.RunMigrations()
	if err != nil {
		return errors.Wrap(err, "Migration")
	}

	//context
	ctx := context.Background()

	// Validator
	v := config.InitValidation()

	repos := repository.NewRepos(logs, db, ctx)

	iter := iteractors.NewIter(logs, repos, ctx)

	delivery := delivery.NewDeliveryService(logs, iter, v)

	// Router
	router := http.NewServeMux()
	router.HandleFunc("/insert", delivery.AddUser)
	router.HandleFunc("/getId/", delivery.GetById)
	router.HandleFunc("/refresh", delivery.RefreshHandler)

	srv := &http.Server{
		Addr:    os.Getenv("REST_PORT"),
		Handler: router,
	}

	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", os.Getenv("REST_PORT"))

	// start HTTP server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Program ended!")

	return nil
}
