package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/longkey1/lnkr/internal/lnkr"
	"github.com/spf13/cobra"
)

var (
	remoteDir        string
	withCreateRemote bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project",
	Long: `Initialize the project by creating necessary configuration files and setting up git exclusions.

This command will:
- Create .lnkr.toml configuration file if it doesn't exist
- Add .lnkr.toml to .git/info/exclude to prevent it from being tracked`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to get current directory: %v\n", err)
			os.Exit(1)
		}

		// Get the number of depth to go up from environment variable, default to DefaultRemoteDepth
		depthStr := os.Getenv("LNK_REMOTE_DEPTH")
		depth := lnkr.DefaultRemoteDepth // default value
		if depthStr != "" {
			if parsedDepth, err := strconv.Atoi(depthStr); err == nil && parsedDepth >= 0 {
				depth = parsedDepth
			}
		}

		// Get base directory for remote
		baseDir := os.Getenv("LNK_REMOTE_ROOT")
		if baseDir == "" {
			// Default to $HOME/.config/lnkr
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to get home directory: %v\n", err)
				os.Exit(1)
			}
			baseDir = filepath.Join(homeDir, ".config", "lnkr")
		}

		// Get remote directory from flag or default
		if remoteDir == "" {
			// Use lnk package function to get default remote path
			remoteDir = lnkr.GetDefaultRemotePath(currentDir, baseDir, depth)
		} else {
			// If remoteDir is specified, make it absolute path based on baseDir
			if !filepath.IsAbs(remoteDir) {
				remoteDir = filepath.Join(baseDir, remoteDir)
			}
		}

		if err := lnkr.Init(remoteDir, withCreateRemote); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&remoteDir, "remote", "r", "", "Remote directory to save in .lnkr.toml (if not specified, uses LNK_REMOTE_ROOT/project-name or parent-dir/current-dir based on LNK_REMOTE_DEPTH)")
	initCmd.Flags().BoolVar(&withCreateRemote, "with-create-remote", false, "Create remote directory if it does not exist")
}
