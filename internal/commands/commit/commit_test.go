package commit

import (
	"fmt"
	"got_it/internal/commands/config"
	init_ "got_it/internal/commands/init"
	"os"
	"strings"
	"testing"
)

// Test GetUserAndEmail function
func TestGetUserAndEmail(t *testing.T) {
	// Create a new Commit instance
	commit := NewCommit("test commit")

	shouldBeUser := "testuser"
	shouldBeEmail := "test@example.com"

	// Set Got environment:
	// Create a temporary directory for testing
	tempdir := t.TempDir()
	err := os.Chdir(tempdir)
	fmt.Println(tempdir)
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}
	// Initialize the Got_it repository
	initializeRepo(t)

	// Set the Got_it configuration
	setKeyValue("user.name", shouldBeUser, t)
	setKeyValue("user.email", shouldBeEmail, t)

	commitMetadata, err := commit.RunCommit()
	if err != nil {
		t.Errorf("Error running commit: %v", err)
	}
	if commitMetadata == "" {
		t.Errorf("Commit metadata is empty")
	}

	// Check for lines with prefixes `tree` and `parent`
	if !strings.Contains(commitMetadata, "tree ") {
		t.Errorf("Commit metadata does not contain a 'tree' line")
	}
	if !strings.Contains(commitMetadata, "parent ") {
		t.Errorf("Commit metadata does not contain a 'parent' line")
	}
	// Check if author and committer are as expected without timestamps
	expectedAuthorPrefix := fmt.Sprintf("author %s <%s>", shouldBeUser, shouldBeEmail)
	expectedCommitterPrefix := fmt.Sprintf("committer %s <%s>", shouldBeUser, shouldBeEmail)
	if !strings.Contains(commitMetadata, expectedAuthorPrefix) {
		t.Errorf("Commit metadata does not contain expected author prefix")
		t.Errorf("Commit metadata: %s", commitMetadata)
	}
	if !strings.Contains(commitMetadata, expectedCommitterPrefix) {
		t.Errorf("Commit metadata does not contain expected committer prefix")
	}
	// Check if the commit message is as expected
	if !strings.Contains(commitMetadata, "test commit") {
		t.Errorf("Commit metadata does not contain the expected commit message")
	}
}

// GetConfig is a helper function to get the config for testing
func (co *Commit) GetConfig(t *testing.T) *config.Config {
	t.Helper()
	return co.conf
}

// SetConfig is a helper function to set the config for testing
func (co *Commit) SetConfig(conf *config.Config, t *testing.T) {
	t.Helper()
	co.conf = conf
}

// InitializeRepo is a helper function to initialize the Got_it repository
func initializeRepo(t *testing.T) {
	t.Helper()
	// Initialize the Got_it repository
	i := init_.NewInit()
	i.InitRepo()
}

// setKeyValue is a helper function to set a key-value pair in the Got_it configuration
func setKeyValue(key, value string, t *testing.T) {
	t.Helper()
	config := config.NewConfig()

	err := config.SetConfigKeyValue(key, value)
	if err != nil {
		t.Fatalf("Error setting config key: %v", err)
	}
}
