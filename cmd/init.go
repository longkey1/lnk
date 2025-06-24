package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/longkey1/lnk/internal/lnk"
	"github.com/spf13/cobra"
)

var (
	remoteDir    string
	createRemote bool
)

const defaultRemoteDepth = 2

// getDefaultRemotePath returns the default remote path based on current directory and parent directories
func getDefaultRemotePath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get the number of depth to go up from environment variable, default to defaultRemoteDepth
	depthStr := os.Getenv("LNK_REMOTE_DEPTH")
	depth := defaultRemoteDepth // default value
	if depthStr != "" {
		if parsedDepth, err := strconv.Atoi(depthStr); err == nil && parsedDepth >= 0 {
			depth = parsedDepth
		}
	}

	// Split the current path into components
	pathComponents := strings.Split(currentDir, string(os.PathSeparator))

	// Remove empty components (happens with absolute paths)
	var cleanComponents []string
	for _, component := range pathComponents {
		if component != "" {
			cleanComponents = append(cleanComponents, component)
		}
	}

	// Adjust depth if we don't have enough components
	if len(cleanComponents) < depth {
		depth = len(cleanComponents)
	}

	// Get the components for the remote path
	// depth=1: current directory only
	// depth=2: parent directory + current directory
	// depth=3: grandparent directory + parent directory + current directory
	startIndex := len(cleanComponents) - depth
	if startIndex < 0 {
		startIndex = 0
	}

	remoteComponents := cleanComponents[startIndex:]

	// Return absolute path by joining with root
	return string(os.PathSeparator) + strings.Join(remoteComponents, string(os.PathSeparator)), nil
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project",
	Long: `Initialize the project by creating necessary configuration files and setting up git exclusions.

This command will:
- Create .lnk.toml configuration file if it doesn't exist
- Add .lnk.toml to .git/info/exclude to prevent it from being tracked`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get remote directory from flag, environment variable, or default
		if remoteDir == "" {
			// First try LNK_REMOTE_ROOT
			remoteRoot := os.Getenv("LNK_REMOTE_ROOT")
			if remoteRoot != "" {
				// Get current directory name as project name
				currentDir, err := os.Getwd()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: failed to get current directory: %v\n", err)
					os.Exit(1)
				}
				projectName := filepath.Base(currentDir)
				remoteDir = filepath.Join(remoteRoot, projectName)
			} else {
				// Use default path based on current directory structure
				defaultPath, err := getDefaultRemotePath()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				remoteDir = defaultPath
			}
		}

		if err := lnk.Init(remoteDir, createRemote); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&remoteDir, "remote", "r", "", "Remote directory to save in .lnk.toml (if not specified, uses LNK_REMOTE_ROOT/project-name or parent-dir/current-dir based on LNK_REMOTE_DEPTH)")
	initCmd.Flags().BoolVar(&createRemote, "create-remote", false, "Create remote directory if it does not exist")
}
