package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/lnkr/internal/lnk"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up files created by init command",
	Long: `Clean up files and changes made by the init command.

This command will:
- Remove .lnkr.toml configuration file if it exists
- Remove .lnkr.toml entry from .git/info/exclude`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := lnk.Clean(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
