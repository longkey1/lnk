package lnk

import (
	"fmt"
	"os"
	"path/filepath"
)

func Unlink() error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(config.Links) == 0 {
		fmt.Printf("No links found in %s\n", ConfigFileName)
		return nil
	}

	// Use source directory as base for resolving link paths
	baseDir := config.Source

	for _, link := range config.Links {
		if err := removeLinkWithBase(link, baseDir); err != nil {
			fmt.Printf("Error removing link for %s: %v\n", link.Path, err)
			continue
		}
	}

	fmt.Println("Link removal completed.")
	return nil
}

func removeLinkWithBase(link Link, baseDir string) error {
	// Resolve absolute path for link
	linkAbs := filepath.Join(baseDir, link.Path)

	if _, err := os.Stat(linkAbs); os.IsNotExist(err) {
		fmt.Printf("Path does not exist, skipping: %s\n", linkAbs)
		return nil
	}

	switch link.Type {
	case LinkTypeHard:
		info, err := os.Stat(linkAbs)
		if err != nil {
			return fmt.Errorf("failed to stat path: %w", err)
		}

		if info.IsDir() {
			if err := os.RemoveAll(linkAbs); err != nil {
				return fmt.Errorf("failed to remove directory: %w", err)
			}
			fmt.Printf("Removed directory: %s\n", linkAbs)
		} else {
			if err := os.Remove(linkAbs); err != nil {
				return fmt.Errorf("failed to remove hard link: %w", err)
			}
			fmt.Printf("Removed hard link: %s\n", linkAbs)
		}
	case LinkTypeSymbolic:
		if err := os.Remove(linkAbs); err != nil {
			return fmt.Errorf("failed to remove symbolic link: %w", err)
		}
		fmt.Printf("Removed symbolic link: %s\n", linkAbs)
	default:
		return fmt.Errorf("unknown link type: %s", link.Type)
	}

	return nil
}
