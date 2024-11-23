package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version     string = "v0.0.0-0"
	showVersion bool
	gotDir      string = ".got"
	maxdepth    int    = 3
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.got_it.yaml)")

	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Show the version of got_it")

}

func showVersionInfo() {
	if showVersion {
		fmt.Println("got_it version: ", version)
		os.Exit(0)
	}
}
