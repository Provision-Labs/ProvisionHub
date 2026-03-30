package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// PluginManifest represents the conf.json structure
type PluginManifest struct {
	ID          string      `json:"id"`
	Creator     string      `json:"creator"`
	Contact     string      `json:"contact"`
	Type        string      `json:"type"`
	Language    string      `json:"language"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	LoadMode    string      `json:"loadMode"`
	Runtime     Runtime     `json:"runtime,omitempty"`
	Behavior    Behavior    `json:"behavior,omitempty"`
	Version     string      `json:"version"`
	Repository  Repository  `json:"repository,omitempty"`
}

// Runtime contains start and transport config
type Runtime struct {
	Start     StartConfig `json:"start,omitempty"`
	Transport Transport   `json:"transport,omitempty"`
}

// StartConfig represents how to start the plugin
type StartConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// Transport represents gRPC transport config
type Transport struct {
	Protocol  string `json:"protocol"`
	Address   string `json:"address"`
	Insecure  bool   `json:"insecure"`
	TimeoutMs int    `json:"timeoutMs"`
}

// Behavior represents plugin behavior settings
type Behavior struct {
	RestartPolicy string      `json:"restartPolicy,omitempty"`
	MaxRestarts   int         `json:"maxRestarts,omitempty"`
	HealthProbe   HealthProbe `json:"healthProbe,omitempty"`
}

// HealthProbe represents health check configuration
type HealthProbe struct {
	RPC        string `json:"rpc"`
	IntervalMs int    `json:"intervalMs"`
	TimeoutMs  int    `json:"timeoutMs"`
}

// Repository represents the plugin repository
type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// PluginsOutput represents the final plugins.json structure
type PluginsOutput struct {
	Version string          `json:"version"`
	Plugins []PluginOutput  `json:"plugins"`
}

// PluginOutput represents a plugin in the output format
type PluginOutput struct {
	ID          string                 `json:"id"`
	Enabled     bool                   `json:"enabled"`
	Type        string                 `json:"type"`
	Language    string                 `json:"language"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	LoadMode    string                 `json:"loadMode"`
	WorkingDir  string                 `json:"workingDir"`
	ConfigFile  string                 `json:"configFile"`
	Start       StartOutput            `json:"start"`
	Transport   Transport              `json:"transport"`
	Metadata    MetadataOutput         `json:"metadata"`
}

// StartOutput represents start config in output
type StartOutput struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

// MetadataOutput represents metadata in output
type MetadataOutput struct {
	Creator    string     `json:"creator"`
	Contact    string     `json:"contact"`
	Version    string     `json:"version"`
	Repository Repository `json:"repository,omitempty"`
}

// scanPluginsDir scans the plugins directory and reads all conf.json files
func scanPluginsDir(baseDir string) ([]PluginManifest, error) {
	var manifests []PluginManifest

	// List type directories (e.g., auth, provisioning)
	typeEntries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read plugins directory: %w", err)
	}

	for _, typeEntry := range typeEntries {
		if !typeEntry.IsDir() {
			continue
		}

		typeDir := filepath.Join(baseDir, typeEntry.Name())

		// List plugin directories within type (e.g., keycloak-go)
		pluginEntries, err := os.ReadDir(typeDir)
		if err != nil {
			continue
		}

		for _, pluginEntry := range pluginEntries {
			if !pluginEntry.IsDir() {
				continue
			}

			pluginDir := filepath.Join(typeDir, pluginEntry.Name())
			confPath := filepath.Join(pluginDir, "conf.json")

			// Read and parse conf.json
			data, err := os.ReadFile(confPath)
			if err != nil {
				// conf.json not found or unreadable, skip
				continue
			}

			var manifest PluginManifest
			if err := json.Unmarshal(data, &manifest); err != nil {
				log.Printf("Warning: failed to parse %s: %v", confPath, err)
				continue
			}

			manifests = append(manifests, manifest)
		}
	}

	return manifests, nil
}

// generatePluginsJSON converts manifests to plugins.json format
func generatePluginsJSON(manifests []PluginManifest, pluginsDir string) PluginsOutput {
	output := PluginsOutput{
		Version: "1.0",
		Plugins: make([]PluginOutput, len(manifests)),
	}

	for i, manifest := range manifests {
		workingDir := filepath.Join("plugins", manifest.Type, manifest.ID)

		pluginOut := PluginOutput{
			ID:          manifest.ID,
			Enabled:     true, // default to enabled
			Type:        manifest.Type,
			Language:    manifest.Language,
			Title:       manifest.Title,
			Description: manifest.Description,
			LoadMode:    manifest.LoadMode,
			WorkingDir:  workingDir,
			ConfigFile:  filepath.Join(workingDir, "conf.json"),
			Start: StartOutput{
				Command: manifest.Runtime.Start.Command,
				Args:    manifest.Runtime.Start.Args,
				Env: map[string]string{
					"AUTH_PLUGIN_LISTEN_ADDR": manifest.Runtime.Transport.Address,
					"AUTH_PLUGIN_PROVIDER":    strings.ToLower(manifest.Type),
					"AUTH_PLUGIN_VERSION":     manifest.Version,
				},
			},
			Transport: manifest.Runtime.Transport,
			Metadata: MetadataOutput{
				Creator:    manifest.Creator,
				Contact:    manifest.Contact,
				Version:    manifest.Version,
				Repository: manifest.Repository,
			},
		}

		output.Plugins[i] = pluginOut
	}

	return output
}

func generatePlugins(pluginsDir, outputFile, rootDir string) error {
	// Resolve root directory
	if rootDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get current directory: %w", err)
		}
		rootDir = cwd
	}

	// Resolve plugins directory
	if pluginsDir == "" {
		pluginsDir = filepath.Join(rootDir, "apps/control-plane/plugins")
	} else if !filepath.IsAbs(pluginsDir) {
		pluginsDir = filepath.Join(rootDir, pluginsDir)
	}

	// Resolve output file
	if outputFile == "" {
		outputFile = filepath.Join(rootDir, "apps/control-plane/plugins.json")
	} else if !filepath.IsAbs(outputFile) {
		outputFile = filepath.Join(rootDir, outputFile)
	}

	log.Printf("Scanning plugins in: %s", pluginsDir)

	// Check if plugins directory exists
	if _, err := os.Stat(pluginsDir); err != nil {
		return fmt.Errorf("plugins directory not found: %s", pluginsDir)
	}

	// Scan plugins
	manifests, err := scanPluginsDir(pluginsDir)
	if err != nil {
		return fmt.Errorf("failed to scan plugins: %w", err)
	}

	if len(manifests) == 0 {
		return fmt.Errorf("no plugins found in %s", pluginsDir)
	}

	log.Printf("Found %d plugin(s)", len(manifests))

	// Generate JSON
	pluginsOutput := generatePluginsJSON(manifests, pluginsDir)

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(pluginsOutput, "", "   ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write output
	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	log.Printf("✓ Generated %s", outputFile)
	for _, plugin := range pluginsOutput.Plugins {
		log.Printf("  - %s (%s@%s) - %s", plugin.ID, plugin.Type, plugin.Metadata.Version, plugin.Title)
	}

	return nil
}
