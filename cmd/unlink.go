package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/lnkr/internal/lnk"
	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Remove links based on .lnkr.toml configuration",
	Long:  `Remove hard links, symbolic links, or directories based on the .lnkr.toml configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := lnk.Unlink(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)
}
