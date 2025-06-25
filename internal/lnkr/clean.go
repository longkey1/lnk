package lnkr

import (
	"fmt"
	"os"
	"strings"
)

// Clean performs the cleanup tasks
func Clean() error {
	// Remove .lnkr.toml file if it exists
	if err := removeLnkToml(); err != nil {
		return fmt.Errorf("failed to remove %s: %w", ConfigFileName, err)
	}

	// Remove .lnkr.toml from .git/info/exclude
	if err := removeFromGitExclude(); err != nil {
		return fmt.Errorf("failed to remove from %s: %w", GitExcludePath, err)
	}

	fmt.Println("Cleanup completed successfully!")
	return nil
}

// removeLnkToml removes the .lnkr.toml file if it exists
func removeLnkToml() error {
	filename := ConfigFileName

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("%s does not exist\n", filename)
		return nil
	}

	// Remove file
	if err := os.Remove(filename); err != nil {
		return err
	}

	fmt.Printf("Removed %s\n", filename)
	return nil
}

// removeFromGitExclude removes .lnkr.toml from .git/info/exclude
func removeFromGitExclude() error {
	// Load config to get git exclude path
	config, err := loadConfig()
	if err != nil {
		// If config doesn't exist, use default path
		return removeFromGitExcludeWithPath(GitExcludePath, ConfigFileName)
	}

	return removeFromGitExcludeWithPath(config.GetGitExcludePath(), ConfigFileName)
}

// removeFromGitExcludeWithPath removes entries from a specific git exclude file
func removeFromGitExcludeWithPath(excludePath, entry string) error {
	// Check if exclude file exists
	if _, err := os.Stat(excludePath); os.IsNotExist(err) {
		fmt.Printf("%s does not exist\n", excludePath)
		return nil
	}

	// Read existing content
	content, err := os.ReadFile(excludePath)
	if err != nil {
		return err
	}

	// Split content into lines
	lines := strings.Split(string(content), "\n")

	// Check if entry exists
	entryExists := false
	for _, line := range lines {
		if strings.TrimSpace(line) == entry {
			entryExists = true
			break
		}
	}

	if !entryExists {
		fmt.Printf("%s does not exist in %s\n", entry, excludePath)
		return nil
	}

	// Filter out the entry
	var newLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != entry {
			newLines = append(newLines, line)
		}
	}

	// Write back the filtered content
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(excludePath, []byte(newContent), 0644); err != nil {
		return err
	}

	fmt.Printf("Removed %s from %s\n", entry, excludePath)
	return nil
}
