package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/lnkr/internal/lnkr"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of links in .lnkr.toml configuration",
	Long:  `Show the status of all links defined in the .lnkr.toml configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := lnkr.Status(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
