package cmd

import (
	"fmt"
	"got_it/internal/config"
	"os"

	"github.com/spf13/cobra"
)

var (
	c           *config.Config
	version     string = "v0.0.0-0"
	showVersion bool
)

var rootCmd = &cobra.Command{
	Use:   "got",
	Short: "",
	Long:  ``,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		showVersionInfo()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show the version of got_it")
}

func showVersionInfo() {
	if showVersion {
		fmt.Println("got_it version: ", version)
		os.Exit(0)
	}
}
