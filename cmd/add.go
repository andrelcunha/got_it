package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [files]",
	Short: "Add files to the staging area",
	Long:  `Add files to the staging area for the next commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Adding files:")
		addFiles(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

// addFiles adds files to the staging area
func addFiles(files []string) {
	currentDepth := 0
	// Ensure the .got directory exists

	if _, err := os.Stat(gotDir); os.IsNotExist(err) {
		fmt.Println("Not a Got_it repository. Run 'got init' first.")
		return
	}

	// Add files to the staging area
	for _, file := range files {
		// Get file information
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			continue
		}

		if fileInfo.IsDir() {
			filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					stageFile(path)
				} else {
					if maxdepth > -1 {
						if currentDepth >= maxdepth+1 {
							return filepath.SkipDir
						} else {
							currentDepth++
						}
					}
				}
				return nil
			})
		} else {
			stageFile(file)
		}
	}
}

func stageFile(file string) {
	fmt.Println("Staging file:", file)
}
