package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/plugins"
)

type report struct {
	PersistentExpected []string
	PersistentStarted  []string
	OnDemandExpected   []string
	OnDemandStarted    []string
}

func main() {
	registryPath := flag.String("registry", envOrDefault("PLUGINS_REGISTRY_PATH", "plugins.json"), "path to plugins registry JSON")
	startupWait := flag.Duration("startup-wait", 1200*time.Millisecond, "wait time after plugin process starts")
	flag.Parse()

	absRegistryPath, err := filepath.Abs(*registryPath)
	if err != nil {
		fatalf("resolve registry path: %v", err)
	}

	reg, err := loadRegistry(absRegistryPath)
	if err != nil {
		fatalf("load registry: %v", err)
	}

	mgr, err := plugins.LoadManager(absRegistryPath)
	if err != nil {
		fatalf("load manager: %v", err)
	}

	r := report{}
	persistentSet := map[string]struct{}{}
	onDemandSet := map[string]struct{}{}

	for _, p := range reg.Plugins {
		if !p.Enabled {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(p.LoadMode)) {
		case "persistent":
			r.PersistentExpected = append(r.PersistentExpected, p.ID)
			persistentSet[p.ID] = struct{}{}
		case "on-demand":
			r.OnDemandExpected = append(r.OnDemandExpected, p.ID)
			onDemandSet[p.ID] = struct{}{}
		default:
			fatalf("plugin %s has invalid loadMode %q (expected persistent or on-demand)", p.ID, p.LoadMode)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	persistentRunning, err := mgr.StartPersistent(ctx)
	if err != nil {
		fatalf("start persistent plugins: %v", err)
	}
	defer stopPlugins(persistentRunning)

	for _, rp := range persistentRunning {
		if err := ensureProcessAlive(rp, *startupWait); err != nil {
			fatalf("persistent plugin %s failed startup check: %v", rp.Config.ID, err)
		}
		r.PersistentStarted = append(r.PersistentStarted, rp.Config.ID)
	}

	sort.Strings(r.PersistentExpected)
	sort.Strings(r.PersistentStarted)
	if !sameIDs(r.PersistentExpected, r.PersistentStarted) {
		fatalf("persistent mismatch: expected=%v started=%v", r.PersistentExpected, r.PersistentStarted)
	}

	for _, id := range r.OnDemandExpected {
		rp, startErr := mgr.StartOnDemand(ctx, id)
		if startErr != nil {
			fatalf("on-demand plugin %s failed start: %v", id, startErr)
		}
		if err := ensureProcessAlive(rp, *startupWait); err != nil {
			fatalf("on-demand plugin %s failed startup check: %v", id, err)
		}
		r.OnDemandStarted = append(r.OnDemandStarted, id)
		stopPlugins([]plugins.RunningPlugin{rp})
	}

	sort.Strings(r.OnDemandExpected)
	sort.Strings(r.OnDemandStarted)
	if !sameIDs(r.OnDemandExpected, r.OnDemandStarted) {
		fatalf("on-demand mismatch: expected=%v started=%v", r.OnDemandExpected, r.OnDemandStarted)
	}

	fmt.Println("[plugin-lifecycle] SUCCESS")
	fmt.Printf("[plugin-lifecycle] persistent expected=%v started=%v\n", r.PersistentExpected, r.PersistentStarted)
	fmt.Printf("[plugin-lifecycle] on-demand expected=%v started=%v\n", r.OnDemandExpected, r.OnDemandStarted)
}

func envOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func loadRegistry(path string) (plugins.Registry, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return plugins.Registry{}, fmt.Errorf("read %s: %w", path, err)
	}

	var reg plugins.Registry
	if err := json.Unmarshal(content, &reg); err != nil {
		return plugins.Registry{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return reg, nil
}

func ensureProcessAlive(rp plugins.RunningPlugin, wait time.Duration) error {
	time.Sleep(wait)
	if rp.Cmd == nil || rp.Cmd.Process == nil {
		return fmt.Errorf("process handle is nil")
	}
	if rp.Cmd.ProcessState != nil && rp.Cmd.ProcessState.Exited() {
		return fmt.Errorf("process exited early")
	}
	if err := rp.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
		return fmt.Errorf("process not alive: %w", err)
	}
	if err := waitTransportReady(rp.Config.Transport, wait); err != nil {
		return fmt.Errorf("transport not ready: %w", err)
	}
	return nil
}

func waitTransportReady(t plugins.PluginTransport, maxWait time.Duration) error {
	address := strings.TrimSpace(t.Address)
	if address == "" {
		return nil
	}

	deadline := time.Now().Add(maxWait)
	for {
		conn, err := net.DialTimeout("tcp", address, 350*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("endpoint %s unreachable: %w", address, err)
		}
		time.Sleep(150 * time.Millisecond)
	}
}

func stopPlugins(running []plugins.RunningPlugin) {
	for _, rp := range running {
		if rp.Cmd == nil || rp.Cmd.Process == nil {
			continue
		}
		_ = rp.Cmd.Process.Kill()
		_ = rp.Cmd.Wait()
	}
}

func sameIDs(expected, got []string) bool {
	if len(expected) != len(got) {
		return false
	}
	for i := range expected {
		if expected[i] != got[i] {
			return false
		}
	}
	return true
}

func fatalf(format string, args ...any) {
	fmt.Printf("[plugin-lifecycle] ERROR: "+format+"\n", args...)
	os.Exit(1)
}
