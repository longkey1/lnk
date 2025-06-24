package lnk

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

// File represents a single file entry in the configuration
type Link struct {
	Path string `toml:"path"`
	Type string `toml:"type"`
}

// Config represents the .lnk.toml configuration structure
type Config struct {
	Links []Link `toml:"links"`
}

// Add adds a file to the project configuration
func Add(path string, recursive bool, linkType string) error {
	// Validate link type
	if linkType != "hard" && linkType != "symbolic" {
		return fmt.Errorf("invalid link type: %s. Must be 'hard' or 'symbolic'", linkType)
	}

	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	existing := make(map[string]struct{})
	for _, link := range config.Links {
		existing[link.Path] = struct{}{}
	}

	var targets []string

	if recursive && fi.IsDir() {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if _, ok := existing[p]; !ok {
					targets = append(targets, p)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		if _, ok := existing[path]; !ok {
			targets = append(targets, path)
		}
	}

	if len(targets) == 0 {
		fmt.Println("No new paths to add.")
		return nil
	}

	for _, t := range targets {
		config.Links = append(config.Links, Link{Path: t, Type: linkType})
		fmt.Printf("Added link: %s (type: %s)\n", t, linkType)
	}

	// pathで昇順ソート
	sort.Slice(config.Links, func(i, j int) bool {
		return config.Links[i].Path < config.Links[j].Path
	})

	if err := saveConfig(config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

// loadConfig loads the configuration from .lnk.toml
func loadConfig() (*Config, error) {
	filename := ".lnk.toml"
	config := &Config{}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create empty config if file doesn't exist
		return config, nil
	}

	// Read and parse existing file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(content) > 0 {
		if _, err := toml.Decode(string(content), config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// saveConfig saves the configuration to .lnk.toml
func saveConfig(config *Config) error {
	filename := ".lnk.toml"

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode and write configuration
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}
