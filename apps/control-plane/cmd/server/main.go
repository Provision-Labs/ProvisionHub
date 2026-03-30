package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/auth"
	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/config"
	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/middleware"
	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/plugins"
)

func main() {
	cfg := config.LoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pluginManager, err := plugins.LoadManager(cfg.PluginsRegistryPath)
	if err != nil {
		log.Fatalf("Failed to load plugins registry: %v", err)
	}

	runningPlugins, err := pluginManager.StartPersistent(ctx)
	if err != nil {
		log.Fatalf("Failed to start persistent plugins: %v", err)
	}
	for _, plugin := range runningPlugins {
		log.Printf("Started plugin %s (pid=%d)", plugin.Config.ID, plugin.Cmd.Process.Pid)
	}

	authPlugin, err := pluginManager.ResolveAuthPlugin()
	if err != nil {
		log.Fatalf("Failed to resolve auth plugin: %v", err)
	}
	if authPlugin.LoadMode == "on-demand" {
		runningPlugin, startErr := pluginManager.StartOnDemand(ctx, authPlugin.ID)
		if startErr != nil {
			log.Fatalf("Failed to start on-demand auth plugin: %v", startErr)
		}
		log.Printf("Started on-demand plugin %s (pid=%d)", runningPlugin.Config.ID, runningPlugin.Cmd.Process.Pid)
	}

	authAddress := authPlugin.Transport.Address
	if authAddress == "" {
		authAddress = cfg.AuthPluginAddr
	}
	authInsecure := authPlugin.Transport.Insecure
	authTimeoutMs := authPlugin.Transport.Timeout
	if authTimeoutMs <= 0 {
		authTimeoutMs = cfg.AuthPluginTimeoutMs
	}

	authProvider, err := auth.NewGRPCProvider(authAddress, time.Duration(authTimeoutMs)*time.Millisecond, authInsecure)
	if err != nil {
		log.Fatalf("Failed to connect auth plugin: %v", err)
	}
	defer func() {
		if closeErr := authProvider.Close(); closeErr != nil {
			log.Printf("Failed to close auth plugin connection: %v", closeErr)
		}
	}()

	auth.SetProvider(authProvider)

	mux := http.NewServeMux()

	// Register features endpoints
	auth.RegisterRoutes(mux)

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
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
