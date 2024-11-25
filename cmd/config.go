package cmd

import (
	"fmt"
	"got_it/internal/config"
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config <key> <value>",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if !NewInit().IsInitialized() {
			return
		}
		runConfig(cmd, args)
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

func runConfig(cmd *cobra.Command, args []string) {
	len := len(args)
	switch len {
	case 0:
		cmd.Help()
	case 1:
		getConfig(args[0])
	case 2:
		setConfig(args[0], args[1])
	default:
		cmd.Help()
	}
	return
}

func setConfig(key, value string) {
	if !config.IsValidKey(key) {
		fmt.Print(invalidateKeyMessage(key))
		return
	} else {
		if err := c.SetConfigKeyValue(key, value); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Config set successfully.")
	}

	return
}

func invalidateKeyMessage(key string) string {
	message := []string{}
	message = append(message, "Error: "+key+" is not a valid config key\n")
	message = append(message, "Valid keys are:\n")

	for _, k := range config.GetAcceptedKeys() {
		message = append(message, "  "+k+"\n")
	}
	return strings.Join(message, "")
}

func getConfig(key string) {
	if !config.IsValidKey(key) {
		fmt.Print(invalidateKeyMessage(key))
		return
	} else {
		value, err := c.GetConfigKeyValue(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(value)
	}
	return
}
