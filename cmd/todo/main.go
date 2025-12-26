package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	// "os/user"
	"syscall"
	"time"

	"github.com/srmty09/Todo-App/internal/config"
	"github.com/srmty09/Todo-App/internal/http/handlers/tasks"
	"github.com/srmty09/Todo-App/internal/http/handlers/users"
	"github.com/srmty09/Todo-App/internal/storage/sqlite"
	"github.com/srmty09/Todo-App/internal/utils/response"
)

func test()http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		response.WriteJson(w,http.StatusAccepted,"welcome to my todo app")
	}
}



func main() {
	// Load config
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	
	router := http.NewServeMux()

	router.HandleFunc("/api", test())
	router.HandleFunc("POST /api/user", users.New(storage))
	router.HandleFunc("POST /api/user/{id}/add_task/",tasks.Add(storage))
	router.HandleFunc("GET /api/user/{id}/todo/",tasks.GetTodo(storage))
	router.HandleFunc("PATCH /api/user/{id}/todo/completed/{task_id}",tasks.CompletedTask(storage))
	router.HandleFunc("PATCH /api/user/{id}/todo/incompleted/{task_id}",tasks.IncompletedTask(storage))
	// router.HandleFunc("DELETE /api/user/{id}/todo/{task_id}",tasks.DeleteTask())

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
