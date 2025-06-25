package lnkr

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateLinks(fromRemote bool) error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(config.Links) == 0 {
		fmt.Printf("No links found in %s\n", ConfigFileName)
		return nil
	}

	// Choose base directory based on fromRemote flag
	var baseDir string
	if fromRemote {
		baseDir = config.Remote
	} else {
		baseDir = config.Local
	}

	for _, link := range config.Links {
		if err := createLinkWithBase(link, baseDir, fromRemote, config); err != nil {
			fmt.Printf("Error creating link for %s: %v\n", link.Path, err)
			continue
		}
	}

	fmt.Println("Link creation completed.")
	return nil
}

func createLinkWithBase(link Link, baseDir string, fromRemote bool, config *Config) error {
	// Determine source and target directories based on fromRemote flag
	var sourceDir, targetDir string
	if fromRemote {
		// When fromRemote is true: remote -> local
		sourceDir = config.Remote
		targetDir = config.Local
	} else {
		// When fromRemote is false: local -> remote
		sourceDir = config.Local
		targetDir = config.Remote
	}

	// Resolve absolute paths for source and target
	sourceAbs := filepath.Join(sourceDir, link.Path)
	targetAbs := filepath.Join(targetDir, link.Path)

	// Check if source exists
	sourceInfo, err := os.Stat(sourceAbs)
	if os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", sourceAbs)
	}

	// Create target directory if it doesn't exist
	targetParentDir := filepath.Dir(targetAbs)
	if err := os.MkdirAll(targetParentDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Check if target already exists
	if _, err := os.Stat(targetAbs); err == nil {
		return fmt.Errorf("target already exists: %s", targetAbs)
	}

	switch link.Type {
	case LinkTypeHard:
		if sourceInfo.IsDir() {
			// For directories, create the directory
			if err := os.MkdirAll(targetAbs, sourceInfo.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			fmt.Printf("Created directory: %s\n", targetAbs)
		} else {
			// For files, create hard link
			if err := os.Link(sourceAbs, targetAbs); err != nil {
				return fmt.Errorf("failed to create hard link: %w", err)
			}
			fmt.Printf("Created hard link: %s -> %s\n", sourceAbs, targetAbs)
		}
	case LinkTypeSymbolic:
		// Create symbolic link
		if err := os.Symlink(sourceAbs, targetAbs); err != nil {
			return fmt.Errorf("failed to create symbolic link: %w", err)
		}
		fmt.Printf("Created symbolic link: %s -> %s\n", sourceAbs, targetAbs)
	default:
		return fmt.Errorf("unknown link type: %s", link.Type)
	}

	return nil
}
