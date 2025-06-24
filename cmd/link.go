package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/lnk/internal/lnk"
	"github.com/spf13/cobra"
)

var sourceRemote bool

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create links based on .lnk.toml configuration",
	Long:  `Create hard links, symbolic links, or directories based on the .lnk.toml configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		sourceRemote, _ := cmd.Flags().GetBool("source-remote")
		if err := lnk.CreateLinks(sourceRemote); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
	linkCmd.Flags().Bool("source-remote", false, "Use remote directory as base for link source paths")
}
