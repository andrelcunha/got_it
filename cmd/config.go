package cmd

import (
	"github.com/spf13/cobra"
)

var acceptedKeys = map[string]string{
	"user.name":  "Your name",
	"user.email": "Your email",
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config <key> <value>",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	configCmd.SetHelpFunc(configHelp)
	rootCmd.AddCommand(configCmd)

}

func configHelp(cmd *cobra.Command, args []string) {
	helpMessage := []string{
		"Usage:",
		"  got config [flags]",
		"  got config [flags] <key> <value>",
		"",
		"Flags:",
		"  -h,    Show this help message",
		"",
	}

	for _, line := range helpMessage {
		cmd.Println(line)
	}
}
