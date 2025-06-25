package lnk

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Configuration file name constant
const ConfigFileName = ".lnkr.toml"

// Git exclude file path constant
const GitExcludePath = ".git/info/exclude"

// Link type constants
const (
	LinkTypeHard     = "hard"
	LinkTypeSymbolic = "symbolic"
)

// Default remote depth constant
const DefaultRemoteDepth = 2

type Link struct {
	Path string `toml:"path"`
	Type string `toml:"type"`
}

type Config struct {
	Source string `toml:"source"`
	Remote string `toml:"remote"`
	Links  []Link `toml:"links"`
}

// GetDefaultRemotePath returns the default remote path based on base directory and remote directory
func GetDefaultRemotePath(baseDir, remoteDir string, depth int) string {
	// Split the base directory path into components
	pathComponents := strings.Split(baseDir, string(os.PathSeparator))

	// Remove empty components (happens with absolute paths)
	var cleanComponents []string
	for _, component := range pathComponents {
		if component != "" {
			cleanComponents = append(cleanComponents, component)
		}
	}

	// Adjust depth if we don't have enough components
	if len(cleanComponents) < depth {
		depth = len(cleanComponents)
	}

	// Get the components for the remote path
	// depth=1: current directory only
	// depth=2: parent directory + current directory
	// depth=3: grandparent directory + parent directory + current directory
	startIndex := len(cleanComponents) - depth
	if startIndex < 0 {
		startIndex = 0
	}

	remoteComponents := cleanComponents[startIndex:]
	remotePath := strings.Join(remoteComponents, string(os.PathSeparator))
	return filepath.Join(remoteDir, remotePath)
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
