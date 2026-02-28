package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/config"
	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/middleware"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("/hello-world", handler)

	addr := ":" + cfg.Port
	log.Printf("Starting server on port %s...\n", cfg.Port)

	handlerWithMw := middleware.RequestLogger(mux)

	srv := &http.Server{
		Addr:         addr,
		Handler:      handlerWithMw,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server running on the configured port
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello World")
	if err != nil {
		panic(err)
	}
}
