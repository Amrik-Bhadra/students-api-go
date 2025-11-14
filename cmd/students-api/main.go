// Entry point

package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Amrik-Bhadra/students-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// databse setup

	//setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to students api"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.Addrress,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addrress))
	
	// note: In production, graceful shutdown is needed

	// creating a channel
	done := make(chan os.Signal, 1)

	// notify in channel, if interrup signal obtained from os
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()


	<- done
	
	// here we will write logic to stop server
	slog.Info("Shutting down the server")
	
	// learn more about context
	// set a timer, that within this time if shutdown doesnot happen, then throw error
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err :=server.Shutdown(ctx);err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
