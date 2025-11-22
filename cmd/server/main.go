package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/stpotter16/biodata/internal/handlers"
	"github.com/stpotter16/biodata/internal/store/db"
	"github.com/stpotter16/biodata/internal/store/sqlite"
)

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbPath := "biodata.sqlite"
	log.Printf("Opening database at %v", dbPath)
	db, err := db.New(dbPath)
	if err != nil {
		return err
	}
	store := sqlite.New(db)
	log.Printf("Opened the store: %v", store)

	handler := handlers.NewServer()
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Received termination signal. Shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("Server shutdown error: %w", err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(
		ctx,
		nil,
		nil,
		nil,
		os.Stdout,
		os.Stderr,
	); err != nil {
		log.Fatalf("%s", err)
	}
}
