package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/lnk/internal/lnk"
	"github.com/spf13/cobra"
)

var (
	remoteDir    string
	createRemote bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project",
	Long: `Initialize the project by creating necessary configuration files and setting up git exclusions.

This command will:
- Create .lnk.toml configuration file if it doesn't exist
- Add .lnk.toml to .git/info/exclude to prevent it from being tracked`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get remote directory from environment variable if not specified via flag
		if remoteDir == "" {
			remoteDir = os.Getenv("LNK_REMOTE")
		}

		if err := lnk.Init(remoteDir, createRemote); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&remoteDir, "remote", "r", "", "Remote directory to save in .lnk.toml (can also be set via LNK_REMOTE environment variable)")
	initCmd.Flags().BoolVar(&createRemote, "create-remote", false, "Create remote directory if it does not exist")
}
