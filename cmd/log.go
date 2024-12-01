package cmd

import (
	"got_it/internal/commands/history"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Commit history of a Git repository",
	Long:  `Shows a chronological list of commits, along with detailed information such as commit hashes, authors, timestamps, and commit messages`,
	Run: func(cmd *cobra.Command, args []string) {
		runLog()
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

}

func runLog() {
	history.Execute()
}
