package lnk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Init performs the initialization tasks
func Init(remote string, createRemote bool) error {
	if err := createLnkTomlWithRemote(remote, createRemote); err != nil {
		return fmt.Errorf("failed to create %s: %w", ConfigFileName, err)
	}

	if err := addToGitExclude(); err != nil {
		return fmt.Errorf("failed to add to %s: %w", GitExcludePath, err)
	}

	fmt.Println("Project initialized successfully!")
	return nil
}

// createLnkTomlWithRemote creates the .lnkr.toml file with remote if it doesn't exist
func createLnkTomlWithRemote(remote string, createRemote bool) error {
	filename := ConfigFileName

	// Get current directory as absolute path for source
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Convert remote to absolute path if provided
	if remote != "" {
		if !filepath.IsAbs(remote) {
			remote, err = filepath.Abs(remote)
			if err != nil {
				return fmt.Errorf("failed to convert remote to absolute path: %w", err)
			}
		}
		// remoteがディレクトリであることを保証
		info, err := os.Stat(remote)
		if os.IsNotExist(err) {
			if createRemote {
				if err := os.MkdirAll(remote, 0755); err != nil {
					return fmt.Errorf("failed to create remote directory: %w", err)
				}
			} else {
				return fmt.Errorf("remote directory does not exist: %s", remote)
			}
		} else if err == nil {
			if !info.IsDir() {
				return fmt.Errorf("remote path exists but is not a directory: %s", remote)
			}
		} else {
			return fmt.Errorf("failed to stat remote directory: %w", err)
		}
	}

	if _, err := os.Stat(filename); err == nil {
		// update remote if already exists
		content, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		var config map[string]interface{}
		if len(content) > 0 {
			if _, err := toml.Decode(string(content), &config); err != nil {
				return err
			}
		}
		// Always update source and remote
		config["source"] = currentDir
		if remote != "" {
			config["remote"] = remote
		}
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(config); err != nil {
			return err
		}
		fmt.Printf("Updated source and remote in %s\n", filename)
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	config := map[string]interface{}{
		"source": currentDir,
	}
	if remote != "" {
		config["remote"] = remote
	}
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	fmt.Printf("Created %s\n", filename)
	return nil
}

// addToGitExclude adds .lnkr.toml to .git/info/exclude
func addToGitExclude() error {
	excludePath := GitExcludePath
	excludeDir := filepath.Dir(excludePath)
	entry := ConfigFileName

	// Create .git/info directory if it doesn't exist
	if err := os.MkdirAll(excludeDir, 0755); err != nil {
		return err
	}

	// Read existing content
	content, err := os.ReadFile(excludePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Check if entry already exists
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == entry {
			fmt.Printf("%s already exists in %s\n", entry, excludePath)
			return nil
		}
	}

	// Append entry to file with comment markers
	file, err := os.OpenFile(excludePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add comment markers and entry
	if _, err := file.WriteString("\n"); err != nil {
		return err
	}
	marker := "# lnkr configuration file"
	if _, err := file.WriteString(marker + "\n"); err != nil {
		return err
	}
	if _, err := file.WriteString(entry + "\n"); err != nil {
		return err
	}

	fmt.Printf("Added %s to %s\n", entry, excludePath)
	return nil
}
