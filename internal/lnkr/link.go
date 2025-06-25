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

	var createdLinks []string
	for _, link := range config.Links {
		if err := createLinkWithBase(link, fromRemote, config); err != nil {
			fmt.Printf("Error creating link for %s: %v\n", link.Path, err)
			continue
		}
		// If err is nil, the link was either created successfully or skipped with a warning
		// Add the target path to the list of created links for git exclude
		if fromRemote {
			createdLinks = append(createdLinks, filepath.Join(config.Local, link.Path))
		} else {
			createdLinks = append(createdLinks, filepath.Join(config.Remote, link.Path))
		}
	}

	// Add created links to .git/info/exclude
	if len(createdLinks) > 0 {
		// Convert absolute paths to relative paths for git exclude
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Warning: failed to get current directory: %v\n", err)
		} else {
			var relativePaths []string
			for _, path := range createdLinks {
				if filepath.IsAbs(path) {
					relPath, err := filepath.Rel(currentDir, path)
					if err != nil {
						// If we can't get relative path, use the original path
						relativePaths = append(relativePaths, path)
					} else {
						relativePaths = append(relativePaths, relPath)
					}
				} else {
					relativePaths = append(relativePaths, path)
				}
			}

			if err := addMultipleToGitExclude(relativePaths, "lnkr created links"); err != nil {
				fmt.Printf("Warning: failed to add links to .git/info/exclude: %v\n", err)
			}
		}
	}

	fmt.Println("Link creation completed.")
	return nil
}

func createLinkWithBase(link Link, fromRemote bool, config *Config) error {
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

	// Check if target already exists
	if _, err := os.Stat(targetAbs); err == nil {
		fmt.Printf("Warning: target already exists: %s\n", targetAbs)
		return nil // Skip this link instead of returning error
	}

	switch link.Type {
	case LinkTypeHard:
		if sourceInfo.IsDir() {
			return fmt.Errorf("hard links cannot be created for directories: %s", sourceAbs)
		} else {
			// For files, create hard link
			targetParentDir := filepath.Dir(targetAbs)
			if err := os.MkdirAll(targetParentDir, 0755); err != nil {
				return fmt.Errorf("failed to create target directory: %w", err)
			}
			if err := os.Link(sourceAbs, targetAbs); err != nil {
				return fmt.Errorf("failed to create hard link: %w", err)
			}
			fmt.Printf("Created hard link: %s -> %s\n", sourceAbs, targetAbs)
		}
	case LinkTypeSymbolic:
		// Create symbolic link (works for both files and directories)
		if err := os.Symlink(sourceAbs, targetAbs); err != nil {
			return fmt.Errorf("failed to create symbolic link: %w", err)
		}
		fmt.Printf("Created symbolic link: %s -> %s\n", sourceAbs, targetAbs)
	default:
		return fmt.Errorf("unknown link type: %s", link.Type)
	}

	return nil
}
