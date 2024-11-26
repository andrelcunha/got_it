package cmd

import (
	"fmt"

	"got_it/internal/commands/commit"

	"github.com/spf13/cobra"
)

var (
	allFlagCommit     bool = false
	verboseFlagCommit bool = false
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit [-a] [-F <file> | -m <message>]",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			fmt.Println(arg)
		}

		runCommit(cmd)
		return
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().BoolVarP(&allFlagCommit, "all", "a", false, "add all changes in tracked files to the commit")
	commitCmd.Flags().StringP("file", "F", "", "read commit message from file")
	commitCmd.Flags().StringP("message", "m", "", "commit message ")
	commitCmd.Flags().BoolVarP(&verboseFlagCommit, "verbose", "v", false, "verbose output")

}

func runCommit(cmd *cobra.Command) {
	msg, err := cmd.Flags().GetString("message")
	if err != nil {
		msg = ""
	}

	commit.Execute(msg)
	return
}
