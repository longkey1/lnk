package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/lnk/internal/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lnk",
	Short: "A link helper CLI tool",
	Long: `lnk is a command line tool for managing and working with links.
It provides various utilities for link manipulation, validation, and management.`,
	Version: version.GetVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lnk - A Link helper")
		fmt.Println("Use --help for more information")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.lnk.yaml)")
}
