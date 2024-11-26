package config

import (
	"bufio"
	"fmt"
	"got_it/internal/models"
	"os"
	"path/filepath"
	"strings"
)

var (
	verbose        bool     = false
	essentialFiles []string = []string{
		".gotignore",
	}
)

var acceptedKeys = map[string]string{
	"init.defaultBranch": "main",
	"user.name":          "Your name",
	"user.email":         "user@example.com",
}

const GOT_DIR string = ".got"
const CONFIG_FILE string = "config"
const GOTIGNORE_FILE string = ".gotignore"
const INDEX_FILE string = "index"
const DEFAULT_BRANCH string = "main"

var INDEX_PATH string = filepath.Join(GOT_DIR, INDEX_FILE)

type Config struct {
	GotDir        string
	IndexFile     string
	GotignoreFile string
	// private settings
	defaultBranch string
	userData      models.User
}

type Callback func(
	*bufio.Scanner,
	*bufio.Writer,
	[]string, // {key} or {key,value}
) (string, error)

func NewConfig() *Config {
	userData := &models.User{
		User:  "",
		Email: "",
	}
	return &Config{
		GotDir:        GOT_DIR,
		IndexFile:     INDEX_FILE,
		GotignoreFile: GOTIGNORE_FILE,
		defaultBranch: DEFAULT_BRANCH,
		userData:      *userData,
	}
}

func (c *Config) GetGotDir() string {
	return GOT_DIR
}

func (c *Config) GetDefaultBranch() string {
	return c.defaultBranch
}
func (c *Config) GetUserName() string {
	if c.userData.User == "" {
		name := acceptedKeys["user.name"]
		name, _ = c.GetConfigKeyValue("user.name")
		c.userData.User = name
	}
	return c.userData.User
}

func (c *Config) GetUserEmail() string {
	if c.userData.Email == "" {
		email := acceptedKeys["user.email"]
		email, _ = c.GetConfigKeyValue("user.email")
		c.userData.Email = email
	}
	return c.userData.Email
}

func (c *Config) GetUserData() models.User {
	if c.userData.User == "" || c.userData.Email == "" {
		c.GetUserName()
		c.GetUserEmail()
	}
	return c.userData
}

// SetConfigKeyValue sets the value for the given configuration key in the
// config file. If the key is not found in the map, it returns an error.
func (c *Config) SetConfigKeyValue(key, value string) error {
	if !IsValidKey(key) {
		return fmt.Errorf("Invalid config key: %s", key)
	}

	err := c.writeConfig(key, value)
	if err != nil {
		return fmt.Errorf("Error saving key %s on config file %s: %s", key, CONFIG_FILE, err.Error())
	}
	return nil
}

// GetConfigKeyValue retrieves the value for the given configuration key from the
// acceptedKeys map. If the key is not found in the map, it returns an error.
func (c *Config) GetConfigKeyValue(key string) (string, error) {
	if !IsValidKey(key) {
		return "", fmt.Errorf("Invalid config key: %s", key)
	}

	value, err := c.readConfig(key)
	if err != nil {
		return "", fmt.Errorf("Error reading key %s from config file %s: %s", key, CONFIG_FILE, err.Error())
	}
	return value, nil
}

// GetSectionNameAndKey returns the section name and key from the given key.
func (c *Config) GetSectionAndKey(key string) (string, string) {
	key = strings.TrimSpace(key)
	sectionAndKey := strings.Split(key, ".")
	section := sectionAndKey[0]
	key = strings.TrimPrefix(key, section+".")

	return section, key
}

// Write key-value pairs to the config file
func (c *Config) writeConfig(key, value string) error {
	value = strings.TrimSpace(value)
	// Get section name from the key
	section, key := c.GetSectionAndKey(key)

	// Open the config file for reading
	configPath := filepath.Join(GOT_DIR, CONFIG_FILE)
	// get the absolute path to the config file
	configPath, err := filepath.Abs(configPath)
	file, err := os.OpenFile(configPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("\n Error opening config file %s: %s", CONFIG_FILE, err.Error())
	}
	defer file.Close()

	// Create a temp file for writing
	gotDirPath, _ := filepath.Abs(GOT_DIR)
	tmpFile, err := os.CreateTemp(gotDirPath, CONFIG_FILE+"_*")
	if err != nil {
		return fmt.Errorf("\n Error creating temp file: %s", err.Error())
	}
	defer os.Remove(tmpFile.Name())

	executeCallbackOnSection(section, key, value, file, tmpFile, writeToSection)

	return nil
}

// Read key-value pairs from the config file
func (c *Config) readConfig(key string) (string, error) {
	// Get section name from the key
	section, key := c.GetSectionAndKey(key)
	//check if config file exists

	// Open the config file for reading
	configPath := filepath.Join(GOT_DIR, CONFIG_FILE)
	file, err := os.OpenFile(configPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return "",
			fmt.Errorf("\n Error opening config file %s: %s", CONFIG_FILE, err.Error())
	}
	defer file.Close()

	// Read the config file line by line
	value, err := executeCallbackOnSection(section, key, "", file, nil, readFromSection)
	if err != nil {
		return "", err
	}
	return value, nil
}

// executeCallbackOnSection executes the provided callback function on the specified section of the config file.
// It scans the config file line by line, looking for the section name, and then calls the callback function
// with the scanner, the current line, the key, and the optional value. If the section is found, the callback
// function is executed and its return values are returned. If the section is not found, an error is returned.
func executeCallbackOnSection(section, key, value string, file *os.File, tmpFile *os.File, action Callback) (string, error) {
	// Open the config file for writing
	scanner := bufio.NewScanner(file)
	var writer *bufio.Writer
	if tmpFile != nil {
		writer = bufio.NewWriter(tmpFile)
	}

	sectionFound := false
	flagPreviousLineEmpty := true
	for scanner.Scan() {
		line := scanner.Text()
		if tmpFile != nil && writer != nil {
			flagThisLineisEmpty := strings.TrimSpace(line) == ""
			if !flagThisLineisEmpty {
				writer.WriteString(strings.TrimRight(line, " ") + "\n")
			} else if flagThisLineisEmpty && !flagPreviousLineEmpty {
				writer.WriteString(strings.TrimRight(line, " ") + "\n")
			}
			flagPreviousLineEmpty = flagThisLineisEmpty
		}
		// Check if the line starts with the section name
		if !strings.HasPrefix(line, "["+section+"]") {
			continue
		}
		sectionFound = true
		// scan section
		var err error
		value, err = action(scanner, writer, []string{key, value})
		if err != nil {
			return "", err
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("\n Error reading config file: %s", err.Error())
	}

	// If the section is not found, add it
	if !sectionFound {
		if tmpFile != nil && writer != nil {
			writer.WriteString("\n[" + section + "]\n")
			newline := fmt.Sprintf("    %s = %s", key, value)
			writer.WriteString(newline + "\n")
		}
	}

	// Flush the writer and rename the temp file to the config file
	if tmpFile != nil && writer != nil {
		writer.Flush()
		configPath := filepath.Join(GOT_DIR, CONFIG_FILE)
		// get the absolute path to the config file
		configPath, err := filepath.Abs(configPath)
		if err != nil {
			return "", fmt.Errorf("\n Error getting absolute path to config file: %s", err.Error())
		}
		if err := os.Rename(tmpFile.Name(), configPath); err != nil {
			return "", fmt.Errorf("\n Error renaming temp file: %s\n", err.Error())
		}
	}
	return value, nil
}

// readFromSection reads the value for the given key from a section
func readFromSection(scanner *bufio.Scanner, writer *bufio.Writer, keyValue []string) (string, error) {
	key := keyValue[0]
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key) {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				return value, nil
			}
		}
		if strings.HasPrefix(line, "[") {
			break
		}
	}
	return "", fmt.Errorf("Key %s not found", key)
}

func writeToSection(scanner *bufio.Scanner, writer *bufio.Writer, keyValue []string) (string, error) {
	key := keyValue[0]
	value := keyValue[1]
	newline := fmt.Sprintf("    %s = %s", key, value)

	_, err := writer.WriteString(newline + "\n")
	if err != nil {
		return "", fmt.Errorf("\n Error writing to tmp file : %s", err.Error())
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, key) {
			return newline, nil
		}
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return "", fmt.Errorf("\n Error writing to tmp file : %s", err.Error())
		}

		// Check if the line starts with a new section
		if strings.HasPrefix(line, "[") {
			return newline, nil
		}
	}
	return "", nil
}
