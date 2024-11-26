package commit

import (
	"fmt"
	"got_it/internal/commands/add"
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
	arrangeEnvironment(t, shouldBeUser, shouldBeEmail)

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

// TestReadStagedFiles
func TestReadStagedFiles(t *testing.T) {
	// ARRANGE:
	// Create a new Commit instance
	commit := NewCommit("test commit")
	// Set Got environment:
	addedFiles, err := arrangeEnvironment(t, "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Error setting up environment: %v", err)
	}
	// ACT:
	// Read staged files
	stagedFiles, err := commit.readStagedFiles()
	if err != nil {
		t.Fatalf("Error reading staged files: %v", err)
	}

	// ASSERT:
	// Check if the staged files are as expected
	for _, file := range addedFiles {
		if _, ok := stagedFiles[file]; !ok {
			t.Errorf("File %s is not in the staged files", file)
		}
	}
}

// HELPER FUNCTIONS

func arrangeEnvironment(t *testing.T, shouldBeUser string, shouldBeEmail string) ([]string, error) {
	// Create a temporary directory
	createTempDir(t)
	//
	initializeRepo(t)
	// Initialize the repository
	setUserAndEmail(shouldBeUser, shouldBeEmail, t)

	// Add a file to the repository
	addedFiles, err := addFilesToRepo(t)
	if err != nil {
		t.Fatalf("Error adding files to repository: %v", err)
	}
	return addedFiles, err

}

func createTempDir(t *testing.T) {
	tempdir := t.TempDir()
	err := os.Chdir(tempdir)
	fmt.Println(tempdir)
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}
}

// setUserAndEmail is a helper function to set the user and email for testing
func setUserAndEmail(shouldBeUser, shouldBeEmail string, t *testing.T) {
	setKeyValue("user.name", shouldBeUser, t)
	setKeyValue("user.email", shouldBeEmail, t)
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

// addFilesToRepo is a helper function to add files to the index
// Creates a list of temp files, adds them to the index, and returns the list of files added
func addFilesToRepo(t *testing.T) ([]string, error) {
	t.Helper()
	// create a list to store the files to be added
	var addedFiles [10]string

	// create temporary files
	for i := 0; i < len(addedFiles); i++ {
		tempFile, err := os.CreateTemp(".", "testfile")
		if err != nil {
			return nil, fmt.Errorf("Error creating temporary file: %v", err)
		}
		addedFiles[i] = tempFile.Name()
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())
		// add the file to the index
	}

	add.Execute(addedFiles[:], true)
	return addedFiles[:], nil
}
