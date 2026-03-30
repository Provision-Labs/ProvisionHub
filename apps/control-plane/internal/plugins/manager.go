package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Registry struct {
	Version string         `json:"version"`
	Plugins []PluginConfig `json:"plugins"`
}

type PluginConfig struct {
	ID          string          `json:"id"`
	Enabled     bool            `json:"enabled"`
	Type        string          `json:"type"`
	Language    string          `json:"language"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	LoadMode    string          `json:"loadMode"`
	WorkingDir  string          `json:"workingDir"`
	ConfigFile  string          `json:"configFile"`
	Start       PluginStart     `json:"start"`
	Transport   PluginTransport `json:"transport"`
	Metadata    PluginMetadata  `json:"metadata"`
}

type PluginStart struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

type PluginTransport struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Insecure bool   `json:"insecure"`
	Timeout  int    `json:"timeoutMs"`
}

type PluginMetadata struct {
	Creator    string `json:"creator"`
	Contact    string `json:"contact"`
	Version    string `json:"version"`
	Repository struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"repository"`
}

type RunningPlugin struct {
	Config PluginConfig
	Cmd    *exec.Cmd
}

type Manager struct {
	registry Registry
	baseDir  string
}

func LoadManager(registryPath string) (*Manager, error) {
	if strings.TrimSpace(registryPath) == "" {
		return nil, fmt.Errorf("plugins registry path is empty")
	}

	absPath, err := filepath.Abs(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve plugins registry path: %w", err)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins registry %s: %w", absPath, err)
	}

	var registry Registry
	if err := json.Unmarshal(content, &registry); err != nil {
		return nil, fmt.Errorf("invalid plugins registry %s: %w", absPath, err)
	}

	return &Manager{
		registry: registry,
		baseDir:  filepath.Dir(absPath),
	}, nil
}

func (m *Manager) StartPersistent(ctx context.Context) ([]RunningPlugin, error) {
	running := make([]RunningPlugin, 0)

	for _, p := range m.registry.Plugins {
		if !p.Enabled || strings.ToLower(p.LoadMode) != "persistent" {
			continue
		}

		rp, err := m.startPlugin(ctx, p)
		if err != nil {
			return nil, err
		}
		running = append(running, rp)
	}

	return running, nil
}

func (m *Manager) StartOnDemand(ctx context.Context, pluginID string) (RunningPlugin, error) {
	for _, p := range m.registry.Plugins {
		if p.ID == pluginID {
			if !p.Enabled {
				return RunningPlugin{}, fmt.Errorf("plugin %s is disabled", pluginID)
			}
			if strings.ToLower(p.LoadMode) != "on-demand" {
				return RunningPlugin{}, fmt.Errorf("plugin %s is not configured as on-demand", pluginID)
			}
			return m.startPlugin(ctx, p)
		}
	}

	return RunningPlugin{}, fmt.Errorf("plugin %s not found", pluginID)
}

func (m *Manager) ResolveAuthPlugin() (PluginConfig, error) {
	for _, p := range m.registry.Plugins {
		if p.Enabled && strings.EqualFold(p.Type, "auth") {
			return p, nil
		}
	}

	return PluginConfig{}, fmt.Errorf("no enabled auth plugin found in registry")
}

func (m *Manager) resolveWorkingDir(dir string) string {
	if strings.TrimSpace(dir) == "" {
		return m.baseDir
	}
	if filepath.IsAbs(dir) {
		return dir
	}
	return filepath.Join(m.baseDir, dir)
}

func flattenEnv(envMap map[string]string) []string {
	if len(envMap) == 0 {
		return nil
	}

	out := make([]string, 0, len(envMap))
	for k, v := range envMap {
		out = append(out, k+"="+v)
	}
	return out
}

func (m *Manager) startPlugin(ctx context.Context, p PluginConfig) (RunningPlugin, error) {
	if strings.TrimSpace(p.Start.Command) == "" {
		return RunningPlugin{}, fmt.Errorf("plugin %s start command is empty", p.ID)
	}

	cmd := exec.CommandContext(ctx, p.Start.Command, p.Start.Args...)
	cmd.Dir = m.resolveWorkingDir(p.WorkingDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), flattenEnv(p.Start.Env)...)

	if err := cmd.Start(); err != nil {
		return RunningPlugin{}, fmt.Errorf("failed to start plugin %s: %w", p.ID, err)
	}

	return RunningPlugin{Config: p, Cmd: cmd}, nil
}
