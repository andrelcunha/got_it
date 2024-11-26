package cmd

import (
	init_ "got_it/internal/init"

	"github.com/spf13/cobra"
)

var i init_.Init

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new repository",
	Long:  `Initialize a new Got_it repository`,
	Run: func(cmd *cobra.Command, args []string) {
		i := init_.NewInit()
		i.InitRepo()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func isInitialized() bool {
	i := init_.NewInit()
	return i.IsInitialized()
}
