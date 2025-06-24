package lnk

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Configuration file name constant
const ConfigFileName = ".lnk.toml"

// Git exclude file path constant
const GitExcludePath = ".git/info/exclude"

// Link type constants
const (
	LinkTypeHard     = "hard"
	LinkTypeSymbolic = "symbolic"
)

type Link struct {
	Path string `toml:"path"`
	Type string `toml:"type"`
}

type Config struct {
	Source string `toml:"source"`
	Remote string `toml:"remote"`
	Links  []Link `toml:"links"`
}

func loadConfig() (*Config, error) {
	filename := ConfigFileName
	config := &Config{}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return config, nil
	}

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

func saveConfig(config *Config) error {
	filename := ConfigFileName

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}
