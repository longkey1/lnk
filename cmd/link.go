package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/lnkr/internal/lnk"
	"github.com/spf13/cobra"
)

var fromRemote bool

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create links based on .lnkr.toml configuration",
	Long:  `Create hard links, symbolic links, or directories based on the .lnkr.toml configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fromRemote, _ := cmd.Flags().GetBool("from-remote")
		if err := lnk.CreateLinks(fromRemote); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
	linkCmd.Flags().Bool("from-remote", false, "Use remote directory as base for link source paths")
}
