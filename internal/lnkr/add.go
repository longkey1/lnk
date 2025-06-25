package lnkr

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func Add(path string, recursive bool, linkType string, fromRemote bool) error {
	if linkType != LinkTypeHard && linkType != LinkTypeSymbolic {
		return fmt.Errorf("invalid link type: %s. Must be '%s' or '%s'", linkType, LinkTypeHard, LinkTypeSymbolic)
	}

	// Check if path is absolute
	if filepath.IsAbs(path) {
		return fmt.Errorf("absolute path is not allowed: %s. Please use relative path", path)
	}

	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Determine base directory for relative paths
	var baseDir string
	if fromRemote {
		if config.Remote == "" {
			return fmt.Errorf("remote directory not configured. Run 'lnkr init --remote <path>' first")
		}
		baseDir = config.Remote
	} else {
		if config.Local == "" {
			return fmt.Errorf("local directory not configured. Run 'lnkr init --local <path>' first")
		}
		baseDir = config.Local
	}

	// Build absolute path and check if file exists
	absPath := filepath.Join(baseDir, path)
	fi, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	if recursive && linkType == LinkTypeSymbolic {
		return fmt.Errorf("recursive option cannot be used with symbolic links")
	}

	existing := make(map[string]struct{})
	for _, link := range config.Links {
		existing[link.Path] = struct{}{}
	}

	var targets []string

	// Add a single path to targets
	if err := addPathToTargets(absPath, baseDir, existing, &targets); err != nil {
		return err
	}

	// Handle symbolic link case
	if linkType == LinkTypeSymbolic {
		// Symbolic links can only handle single files/directories (no recursive)
	} else {
		// Handle hard link case
		if fi.IsDir() {
			// Directory with hard links requires recursive option
			if !recursive {
				return fmt.Errorf("recursive option must be set when adding a directory with hard links")
			}

			// Walk directory and add all files
			err := filepath.Walk(absPath, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				// Skip directories, only add files for hard links
				if info.IsDir() {
					return nil
				}
				return addPathToTargets(p, baseDir, existing, &targets)
			})
			if err != nil {
				return fmt.Errorf("failed to walk directory: %w", err)
			}
			if len(targets) == 0 {
				return fmt.Errorf("no files or directories to add under the specified directory")
			}
		} else {
			// Single file with hard links
			if err := addPathToTargets(absPath, baseDir, existing, &targets); err != nil {
				return err
			}
		}
	}

	if len(targets) == 0 {
		fmt.Println("No new paths to add.")
		return nil
	}

	// Add links to config
	for _, t := range targets {
		config.Links = append(config.Links, Link{Path: t, Type: linkType})
		fmt.Printf("Added link: %s (type: %s)\n", t, linkType)
	}

	sort.Slice(config.Links, func(i, j int) bool {
		return config.Links[i].Path < config.Links[j].Path
	})

	if err := saveConfig(config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Add target paths to .git/info/exclude
	if err := addTargetPathsToGitExclude(targets, fromRemote, config); err != nil {
		fmt.Printf("Warning: failed to add target paths to .git/info/exclude: %v\n", err)
	}

	return nil
}

// addPathToTargets adds a single path to the targets slice if it doesn't already exist
func addPathToTargets(absPath, baseDir string, existing map[string]struct{}, targets *[]string) error {
	relPath, err := filepath.Rel(baseDir, absPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}
	if _, ok := existing[relPath]; !ok {
		*targets = append(*targets, relPath)
	}
	return nil
}

// addTargetPathsToGitExclude adds target paths of added links to .git/info/exclude
func addTargetPathsToGitExclude(targetPaths []string, fromRemote bool, config *Config) error {
	// Determine target directory based on fromRemote flag
	var targetDir string
	if fromRemote {
		targetDir = config.Local
	} else {
		targetDir = config.Remote
	}

	// Build absolute target paths
	var absoluteTargetPaths []string
	for _, path := range targetPaths {
		absoluteTargetPaths = append(absoluteTargetPaths, filepath.Join(targetDir, path))
	}

	// Convert absolute paths to relative paths for git exclude
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var relativePaths []string
	for _, path := range absoluteTargetPaths {
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

	// Add to .git/info/exclude
	return addMultipleToGitExclude(relativePaths)
}
