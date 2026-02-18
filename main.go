package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grey/database"
	"github.com/grey/models"
	"github.com/grey/routers"
)

func main() {

	// initialize database
	db := database.InitDb()

	// automigration listen to changes in the model
	migrations := database.Migrations{
		DB: db,
		Models: []interface{}{
			&models.User{},
			&models.Account{},
			&models.Payment{},
			&models.LedgerEntry{},
		},
	}
	database.RunMigrations(migrations)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: routers.NewRouter(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("Server exiting")
}
