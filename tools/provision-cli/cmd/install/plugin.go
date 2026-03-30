package main

import (
	"fmt"
	"log"
	"strings"
)

// installPlugin installs a plugin from the registry
// Format: plugin should be "type/id@version" or just "auth@latest"
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
