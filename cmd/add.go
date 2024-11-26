package cmd

import (
	"fmt"
	"got_it/internal/commands/add"

	"github.com/spf13/cobra"
)

var (
	verboseAdd bool
)

var addCmd = &cobra.Command{
	Use:   "add [flags] <files>...",
	Short: "Add files to the staging area",
	Long:  `Add files to the staging area for the next commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No files specified.")
			cmd.Help()
			return
		}
		runAdd(args, verboseAdd)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&verboseAdd, "verbose", "v", false, "be verbose")
	rootCmd.AddCommand(addCmd)
}

// runAdd adds files to the staging area
func runAdd(files []string, verbose bool) {
	add.Execute(files, verbose)
}
