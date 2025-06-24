package lnk

import (
	"fmt"
	"os"
	"strings"
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
		status := checkLinkStatus(link)
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

func checkLinkStatus(link Link) LinkStatus {
	status := LinkStatus{
		Path: link.Path,
		Type: link.Type,
	}

	info, err := os.Stat(link.Path)
	if os.IsNotExist(err) {
		status.Exists = false
		status.Error = "NOT FOUND"
		return status
	}

	status.Exists = true

	if link.Type == LinkTypeSymbolic {
		if info.Mode()&os.ModeSymlink != 0 {
			status.IsLink = true
		} else {
			status.Error = "Not a symbolic link"
		}
	} else if link.Type == LinkTypeHard {
		if info.IsDir() {
			status.IsLink = true
		} else {
			status.IsLink = true
		}
	}

	return status
}
