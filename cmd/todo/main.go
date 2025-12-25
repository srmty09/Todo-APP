package main

import (
	"context"
	// "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/srmty09/Todo-App/internal/config"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	router := http.NewServeMux()

	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the todo app"))
	})

	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("server starting", slog.String("addr", server.Addr))

	// Channel to listen for OS signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Block until signal is received
	<-done
	slog.Info("shutdown signal received")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", slog.Any("error", err))
	} else {
		slog.Info("server stopped gracefully")
	}
}
