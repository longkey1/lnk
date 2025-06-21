package lnk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Init performs the initialization tasks
func Init() error {
	// Create .lnk.toml file if it doesn't exist
	if err := createLnkToml(); err != nil {
		return fmt.Errorf("failed to create .lnk.toml: %w", err)
	}

	// Add .lnk.toml to .git/info/exclude
	if err := addToGitExclude(); err != nil {
		return fmt.Errorf("failed to add to .git/info/exclude: %w", err)
	}

	fmt.Println("Project initialized successfully!")
	return nil
}

// createLnkToml creates the .lnk.toml file if it doesn't exist
func createLnkToml() error {
	filename := ".lnk.toml"

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("%s already exists\n", filename)
		return nil
	}

	// Create empty file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Printf("Created %s\n", filename)
	return nil
}

// addToGitExclude adds .lnk.toml to .git/info/exclude
func addToGitExclude() error {
	excludePath := ".git/info/exclude"
	excludeDir := filepath.Dir(excludePath)
	entry := ".lnk.toml"

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

	// Append entry to file
	file, err := os.OpenFile(excludePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(entry + "\n"); err != nil {
		return err
	}

	fmt.Printf("Added %s to %s\n", entry, excludePath)
	return nil
}
