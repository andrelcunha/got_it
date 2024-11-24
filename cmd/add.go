package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	// Get the absolute path of the repository root
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("Error getting repository root:", err)
		return
	}

	// Add files to the staging area
	for _, file := range files {
		// Get the absolute path of the file
		absFile, err := filepath.Abs(file)
		if err != nil {
			fmt.Printf("Error getting absolute path of %v\n", err)
			continue
		}

		// Check if the file is within the repository
		if repoRoot != absFile {
			if mached, err := filepath.Match(repoRoot+"/*", absFile); err != nil || !mached {
				fmt.Printf("Error: %s is outside the repository\n", absFile)
				fmt.Println("root:", repoRoot)
				continue
			}
		}

		// Get file information
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			continue
		}

		if fileInfo.IsDir() {
			filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					// absFile, err := filepath.Abs(file)

					// if err != nil {
					// 	fmt.Printf("Error getting absolute path of %s\n", info.Name())
					// }
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
	if ignoreFile(file) {
		fmt.Println("Ignoring file:", file)
		return
	}
	fmt.Println("Staging file:", file)
}

// ignoreFile tests if file matches the ignore patterns on .gotignore
func ignoreFile(file string) bool {
	shallIgnore := false
	// matched := false

	// Read .gotignore file
	ignoreFile, err := os.Open(".gotignore")
	if err != nil {
		return false
	}
	defer ignoreFile.Close()
	scanner := bufio.NewScanner(ignoreFile)

	for scanner.Scan() {
		line := scanner.Text()
		pattern := strings.TrimSpace(strings.Split(line, "#")[0])

		if pattern != "" {
			// Check if file matches '!' + pattern
			if strings.HasPrefix(pattern, "!") {
				pattern = pattern[1:]
				notIgnoreMatched, err := filepath.Match(pattern, file)

				if err != nil {
					// fmt.Println("Error matching pattern:", err)
					continue
				}
				if notIgnoreMatched {
					// fmt.Printf("Found '!' on pattern pattern: %s, matched: %v notIgore%v\n", pattern, matched, notIgnoreMatched)
					shallIgnore = false
					break
				}
			}

			matched, err := filepath.Match(pattern, file)
			if err != nil {
				//fmt.Println("Error matching pattern:", err)
				continue
			}
			if matched {
				// fmt.Printf("Found pattern: %s, matched: %v\n", pattern, matched)
				shallIgnore = true
				continue
			}
		}
	}
	return shallIgnore
}
