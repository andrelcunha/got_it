package config

import (
	"fmt"
	"os"
	"strings"
)

var (
	essentialFiles []string = []string{
		".gotignore",
	}
)

var acceptedKeys = map[string]string{
	"init.defaultBranch": "main",
	"user.name":          "Your name",
	"user.email":         "user@example.com",
}

const gotDir string = ".got"
const configFile string = ".got/config"
const maxdepth int = -1
const gotIgnoreFile string = ".gotignore"

type Config struct {
	maxdepth      int    //= -1
	defaultBranch string //= "main"
	userName      string //= ""
	userEmail     string //= ""
}

func NewConfig() *Config {
	return &Config{
		maxdepth:      -1,
		defaultBranch: "main",
		userName:      "",
		userEmail:     "",
	}
}

func (c *Config) GetGotDir() string {
	return gotDir
}
func (c *Config) GetMaxDepth() int {
	return c.maxdepth
}
func (c *Config) GetDefaultBranch() string {
	return c.defaultBranch
}
func (c *Config) GetUserName() string {
	return c.userName
}
func (c *Config) GetUserEmail() string {
	return c.userEmail
}

// SetConfigKeyValue sets the value for the given configuration key in the
// config file. If the key is not found in the map, it returns an error.
func (c *Config) SetConfigKeyValue(key, value string) error {
	if _, ok := acceptedKeys[key]; !ok {
		return fmt.Errorf("Invalid config key: %s", key)
	}
	//TODO: save key-value pair on config file
	return nil
}

// GetConfigKeyValue retrieves the value for the given configuration key from the
// acceptedKeys map. If the key is not found in the map, it returns an error.
func (c *Config) GetConfigKeyValue(key string) (string, error) {
	value, ok := acceptedKeys[key]
	// if _, ok := acceptedKeys[key]; !ok {
	if !ok {
		return "", fmt.Errorf("Invalid config key: %s", key)
	}
	return value, nil
}

// Write key-value pairs to the config file
func (c *Config) writeConfig(key, value string) error {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	// Get session name from the key, the part before the first dot
	sessionAndKey := strings.Split(key, ".")
	session := sessionAndKey[0]
	// trim the session name from the key
	key = strings.TrimPrefix(key, session+".")

	// Open the config file for writing
	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the session name to the config file
	_, err = fmt.Fprintf(file, "[%s]\n", session)
	if err != nil {
		return err
	}

	// Write the key-value pairs to the config file
	_, err = fmt.Fprintf(file, "%s=%s\n", key, value)
	if err != nil {
		return err
	}

	//TODO
	return nil
}
