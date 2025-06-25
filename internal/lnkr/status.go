package lnkr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type LinkStatus struct {
	Path   string
	Type   string
	Exists bool
	IsLink bool
	Error  string
}

func Status() error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(config.Links) == 0 {
		fmt.Printf("No links found in %s\n", ConfigFileName)
		return nil
	}

	var statuses []LinkStatus
	for _, link := range config.Links {
		status := checkLinkStatus(link, config)
		statuses = append(statuses, status)
	}

	// Calculate max width for each column
	maxPath := len("Path")
	maxType := len("Type")
	maxStatus := len("Status")
	for _, s := range statuses {
		if len(s.Path) > maxPath {
			maxPath = len(s.Path)
		}
		if len(s.Type) > maxType {
			maxType = len(s.Type)
		}
		st := getStatusText(s)
		if len(st) > maxStatus {
			maxStatus = len(st)
		}
	}

	// Print header
	header := fmt.Sprintf("%-*s  %-*s  %-*s", maxPath, "Path", maxType, "Type", maxStatus, "Status")
	sep := strings.Repeat("-", len(header))
	fmt.Println(header)
	fmt.Println(sep)

	// Print each status
	for _, s := range statuses {
		st := getStatusText(s)
		fmt.Printf("%-*s  %-*s  %-*s\n", maxPath, s.Path, maxType, s.Type, maxStatus, st)
	}

	return nil
}

func getStatusText(status LinkStatus) string {
	if !status.Exists {
		return "NOT FOUND"
	}
	if status.Error != "" {
		return "ERROR: " + status.Error
	}
	if status.IsLink {
		return "LINKED"
	}
	return "NOT LINKED"
}

func checkLinkStatus(link Link, config *Config) LinkStatus {
	status := LinkStatus{
		Path: link.Path,
		Type: link.Type,
	}

	// Check if the link path exists
	info, err := os.Stat(link.Path)
	if os.IsNotExist(err) {
		status.Exists = false
		status.Error = "NOT FOUND"
		return status
	}

	status.Exists = true

	switch link.Type {
	case LinkTypeSymbolic:
		// Check if it's actually a symbolic link
		if info.Mode()&os.ModeSymlink == 0 {
			status.Error = "Not a symbolic link"
			return status
		}

		// Get the target of the symbolic link
		target, err := os.Readlink(link.Path)
		if err != nil {
			status.Error = fmt.Sprintf("Cannot read link target: %v", err)
			return status
		}

		// Check if the target exists
		if _, err := os.Stat(target); os.IsNotExist(err) {
			status.Error = "Target does not exist"
			return status
		}

		// Check if the target path is correct (should point to remote location)
		expectedTarget := getExpectedTargetPath(link.Path, config)
		if target != expectedTarget {
			status.Error = fmt.Sprintf("Wrong target: %s (expected: %s)", target, expectedTarget)
			return status
		}

		status.IsLink = true

	case LinkTypeHard:
		// Check if the file exists and is not a directory
		if info.IsDir() {
			status.Error = "Hard links cannot be created for directories"
			return status
		}

		// Get the expected target path
		expectedTarget := getExpectedTargetPath(link.Path, config)
		if expectedTarget == "" {
			status.Error = "Cannot determine expected target path"
			return status
		}

		// Check if the target file exists
		targetInfo, err := os.Stat(expectedTarget)
		if os.IsNotExist(err) {
			status.Error = "Target file does not exist"
			return status
		}
		if err != nil {
			status.Error = fmt.Sprintf("Cannot access target file: %v", err)
			return status
		}

		// Check if both files have the same inode (hard link check)
		if info.Sys() == nil || targetInfo.Sys() == nil {
			status.Error = "Cannot get file system info for inode comparison"
			return status
		}

		// Compare inodes to verify hard link
		linkInode := getInode(info)
		targetInode := getInode(targetInfo)

		if linkInode == 0 || targetInode == 0 {
			status.Error = "Cannot determine inode numbers"
			return status
		}

		if linkInode != targetInode {
			status.Error = "Not a hard link (different inodes)"
			return status
		}

		status.IsLink = true
	}

	return status
}

func getExpectedTargetPath(linkPath string, config *Config) string {
	// Get the relative path from the local directory
	relPath, err := filepath.Rel(config.Local, linkPath)
	if err != nil {
		return ""
	}

	// Construct the expected target path in the remote directory
	return filepath.Join(config.Remote, relPath)
}

func getInode(info os.FileInfo) uint64 {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
}
