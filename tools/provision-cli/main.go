package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	Version = "0.1.0"
)

func init() {
	log.SetPrefix("provision-cli: ")
	log.SetFlags(0)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "generate":
		generateCommand(args)
	case "install":
		installCommand(args)
	case "version":
		fmt.Printf("provision-cli v%s\n", Version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`provision-cli v%s - ProvisionHub plugin manager CLI

USAGE:
  provision-cli <command> [options]

COMMANDS:
  generate    Generate plugins.json from conf.json manifests
  install     Install a plugin from registry.provisionlabs.dev
  version     Show version
  help        Show this help message

EXAMPLES:
  # Generate plugins.json from plugins/ directory
  provision-cli generate --plugins-dir ./plugins --output ./plugins.json

  # Generate with specific project root
  provision-cli generate --root /path/to/ProvisionHub

  # Install plugin from registry
  provision-cli install --plugin auth@latest --registry registry.provisionlabs.dev

RUN 'provision-cli <command> -h' for more information on a command.
`, Version)
}

func generateCommand(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf(`provision-cli generate - Generate plugins.json from conf.json manifests

USAGE:
  provision-cli generate [options]

OPTIONS:
`)
		fs.PrintDefaults()
	}

	pluginsDir := fs.String("plugins-dir", "", "Path to plugins directory (relative to project root)")
	output := fs.String("output", "", "Output plugins.json file path")
	root := fs.String("root", "", "Project root directory (auto-detects if not set)")

	fs.Parse(args)

	if err := generatePlugins(*pluginsDir, *output, *root); err != nil {
		log.Fatalf("Failed to generate plugins.json: %v", err)
	}
}

func installCommand(args []string) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf(`provision-cli install - Install a plugin from registry

USAGE:
  provision-cli install [options]

OPTIONS:
`)
		fs.PrintDefaults()
	}

	plugin := fs.String("plugin", "", "Plugin reference (e.g., auth@latest, provisioning@0.1.0)")
	registry := fs.String("registry", "registry.provisionlabs.dev", "Registry URL")
	dest := fs.String("dest", "", "Destination directory (defaults to ./plugins/<type>/<id>)")

	fs.Parse(args)

	if *plugin == "" {
		log.Fatal("--plugin flag is required")
	}

	parts := strings.Split(*plugin, "@")
	if len(parts) < 2 {
		log.Fatalf("Invalid plugin reference format. Expected: <type>/<id>@<version>, got: %s", *plugin)
	}

	if err := installPlugin(*plugin, *registry, *dest); err != nil {
		log.Fatalf("Failed to install plugin: %v", err)
	}
}

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
	// Internal field to track actual directory path
	_actualPath string `json:"-"`
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
// Supports two directory structures:
// 1. plugins/{type}/{plugin-id}/conf.json
// 2. plugins/{plugin-id}/conf.json (type from conf.json)
func scanPluginsDir(baseDir string) ([]PluginManifest, error) {
	var manifests []PluginManifest

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read plugins directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryPath := filepath.Join(baseDir, entry.Name())
		relPath := entry.Name()

		// Check if this is a plugin directory with conf.json (flat structure)
		confPath := filepath.Join(entryPath, "conf.json")
		if _, err := os.Stat(confPath); err == nil {
			// Found conf.json, read it (flat structure)
			if manifest, err := readPluginManifest(confPath); err == nil {
				// For flat structure, actualPath is just the plugin id
				manifest._actualPath = relPath
				manifests = append(manifests, manifest)
				continue
			}
		}

		// Otherwise, treat as type directory and scan subdirectories (hierarchical)
		typeEntries, err := os.ReadDir(entryPath)
		if err != nil {
			continue
		}

		for _, typeEntry := range typeEntries {
			if !typeEntry.IsDir() {
				continue
			}

			pluginDir := filepath.Join(entryPath, typeEntry.Name())
			confPath := filepath.Join(pluginDir, "conf.json")

			if manifest, err := readPluginManifest(confPath); err == nil {
				// For hierarchical structure, actualPath is type/plugin-id
				manifest._actualPath = filepath.Join(relPath, typeEntry.Name())
				manifests = append(manifests, manifest)
			}
		}
	}

	return manifests, nil
}

// readPluginManifest reads and parses a single plugin manifest file
func readPluginManifest(confPath string) (PluginManifest, error) {
	var manifest PluginManifest

	data, err := os.ReadFile(confPath)
	if err != nil {
		return manifest, fmt.Errorf("cannot read file: %w", err)
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return manifest, fmt.Errorf("cannot parse JSON: %w", err)
	}

	return manifest, nil
}

// generatePluginsJSON converts manifests to plugins.json format
func generatePluginsJSON(manifests []PluginManifest, pluginsDir string) PluginsOutput {
	output := PluginsOutput{
		Version: "1.0",
		Plugins: make([]PluginOutput, len(manifests)),
	}

	for i, manifest := range manifests {
		// Use actual path from scanner (preserves original structure)
		workingDir := filepath.Join("plugins", manifest._actualPath)

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

func installPlugin(plugin, registry, dest string) error {
	// Parse plugin reference
	parts := strings.Split(plugin, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid plugin reference format: %s", plugin)
	}

	pluginRef := parts[0]
	version := parts[1]

	var pluginType, pluginID string
	if strings.Contains(pluginRef, "/") {
		// Format: type/id
		refParts := strings.Split(pluginRef, "/")
		if len(refParts) != 2 {
			return fmt.Errorf("invalid plugin reference format: %s", pluginRef)
		}
		pluginType = refParts[0]
		pluginID = refParts[1]
	} else {
		// Assume type is "auth" if not specified (default for backward compatibility)
		// In future, this could be resolved from registry
		pluginType = "auth"
		pluginID = pluginRef
	}

	log.Printf("Installing plugin: %s/%s@%s from %s", pluginType, pluginID, version, registry)

	// TODO: Implement registry fetch logic
	// 1. Fetch plugin metadata from registry: GET /plugins/{type}/{id}/{version}/manifest.json
	// 2. Download plugin bundle: GET /plugins/{type}/{id}/{version}/bundle.tar.gz
	// 3. Extract to destination: dest or ./plugins/{type}/{id}
	// 4. Run post-install hooks (if any)
	// 5. Generate/update plugins.json

	return fmt.Errorf("not yet implemented: registry fetch from %s for %s/%s@%s", registry, pluginType, pluginID, version)
}
