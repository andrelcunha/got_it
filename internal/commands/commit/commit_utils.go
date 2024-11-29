package commit

import (
	"fmt"
	"os"
	"path/filepath"
)

func separator() string {
	return string(filepath.Separator)
}

// getEnvVarValue returns the value of the given environment variable
func getEnvVarValue(localVar *string, envVarName string) error {
	for _, envVar := range supportedEnvVars {
		if envVar == envVarName {
			value := os.Getenv(envVar)
			if value != "" {
				*localVar = value
				return nil
			}
		}
	}
	return fmt.Errorf("Unsupported environment variable: %s", envVarName)
}
