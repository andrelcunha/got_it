package config

import (
	"os"
	"testing"
)

// Test writeConfig function
func TestWriteConfig(t *testing.T) {
	// Create a temporary config file for testing
	tmpfile, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatal(err)
		defer os.Remove(tmpfile.Name())

		// Create a new Config instance
		config := NewConfig()
		// Write a key-value pair to the config file
		err = config.writeConfig("key", "value")
	}
}
