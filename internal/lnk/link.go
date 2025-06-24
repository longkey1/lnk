package lnk

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateLinks(sourceRemote bool) error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(config.Links) == 0 {
		fmt.Printf("No links found in %s\n", ConfigFileName)
		return nil
	}

	// base directory for resolving link source
	var baseDir string
	if sourceRemote {
		if config.Remote == "" {
			return fmt.Errorf("remote directory not configured. Run 'lnk init --remote <path>' first")
		}
		baseDir = config.Remote
	} else {
		baseDir = config.Source
	}

	for _, link := range config.Links {
		if err := createLinkWithBase(link, baseDir); err != nil {
			fmt.Printf("Error creating link for %s: %v\n", link.Path, err)
			continue
		}
	}

	fmt.Println("Link creation completed.")
	return nil
}

func createLinkWithBase(link Link, baseDir string) error {
	// Resolve absolute path for source
	sourceAbs := filepath.Join(baseDir, link.Path)

	// Check if source exists
	sourceInfo, err := os.Stat(sourceAbs)
	if os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", sourceAbs)
	}

	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(sourceAbs)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	switch link.Type {
	case LinkTypeHard:
		if sourceInfo.IsDir() {
			// For directories, create the directory
			if err := os.MkdirAll(sourceAbs, sourceInfo.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			fmt.Printf("Created directory: %s\n", sourceAbs)
		} else {
			// For files, create hard link
			if err := os.Link(sourceAbs, sourceAbs); err != nil {
				return fmt.Errorf("failed to create hard link: %w", err)
			}
			fmt.Printf("Created hard link: %s\n", sourceAbs)
		}
	case LinkTypeSymbolic:
		// Create symbolic link
		if err := os.Symlink(sourceAbs, sourceAbs); err != nil {
			return fmt.Errorf("failed to create symbolic link: %w", err)
		}
		fmt.Printf("Created symbolic link: %s\n", sourceAbs)
	default:
		return fmt.Errorf("unknown link type: %s", link.Type)
	}

	return nil
}
