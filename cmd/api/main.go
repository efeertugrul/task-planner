package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"todo-planning/cmd/api/server"
	"todo-planning/internal/db"
	"todo-planning/internal/logger"

	"gorm.io/gorm"
)

var (
	database *gorm.DB
)

func init() {
	var err error

	database, err = db.NewConnection()
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to database: %w", err))
		panic(err)
	}
}

func run(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := server.NewServer(8080, database)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", router.Port),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go run(srv)

	logger.Info("Server is running on port: ", router.Port)
	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
