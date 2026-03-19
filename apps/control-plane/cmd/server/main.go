package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/auth"
	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/config"
	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	mux := http.NewServeMux()

	// Register features endpoints
	auth.RegisterRoutes(mux)

	// Set Session Secret with config handler
	auth.Init(cfg)
	middleware.InitStore(cfg.SessionSecret)
	auth.InitStore(cfg.SessionSecret)
	config.ConnectDatabase(cfg)

	addr := ":" + strconv.Itoa(cfg.Port)
	log.Printf("Starting server on port %d...\n", cfg.Port)

	// Handling all middlewares needed
	handler := middleware.Chain(mux,
		middleware.RequestLogger,
		middleware.RequireAuth,
	)

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server running on the configured port
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
