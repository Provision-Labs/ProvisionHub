package main

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetPrefix("provision-cli: ")
	log.SetFlags(0)
}

func main() {
	// This file exists to allow stand-alone execution during development
	// In production, use: go run ../../ generate [options]

	cmd := os.Args[0]

	if len(os.Args) < 2 {
		log.Fatal("usage: provision-cli generate [options]")
	}

	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	pluginsDir := fs.String("plugins-dir", "", "Path to plugins directory")
	output := fs.String("output", "", "Output plugins.json file path")
	root := fs.String("root", "", "Project root directory")

	fs.Parse(os.Args[1:])

	if err := generatePlugins(*pluginsDir, *output, *root); err != nil {
		log.Fatalf("Failed: %v", err)
	}
}
