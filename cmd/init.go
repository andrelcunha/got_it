package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new repository",
	Long:  `Initialize a new Got_it repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing a new repository...")
		initRepo()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initRepo() {
	// Check if the .got directory already exists
	if _, err := os.Stat(gotDir); !os.IsNotExist(err) {
		fmt.Println("Repository already initialized.")
		return
	}

	// Create the .got directory
	if err := os.Mkdir(gotDir, 0755); err != nil {
		fmt.Println("Error creating .got directory:", err)
		return
	}
	// Create the .got/objects directory
	if err := os.Mkdir(gotDir+"/objects", 0755); err != nil {
		fmt.Println("Error creating .got/objects directory:", err)
		return
	}

	fmt.Println("Repository initialized successfully.")

}
