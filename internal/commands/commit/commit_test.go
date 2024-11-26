package commit

import (
	"fmt"
	"got_it/internal/commands/config"
	init_ "got_it/internal/commands/init"
	"os"
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

	// Call the GetUserAndEmail function
	user, email := commit.GetUserAndEmail()

	// Check if the returned values are correct
	if user != shouldBeUser || email != shouldBeEmail {
		t.Errorf("GetUserAndEmail() = (%s, %s), want (%s, %s)", user, email, shouldBeUser, shouldBeEmail)
	}
}

func (co *Commit) GetConfig(t *testing.T) *config.Config {
	t.Helper()
	return co.conf
}

func (co *Commit) SetConfig(conf *config.Config, t *testing.T) {
	t.Helper()
	co.conf = conf
}

func initializeRepo(t *testing.T) {
	t.Helper()
	// Initialize the Got_it repository
	i := init_.NewInit()
	i.InitRepo()
}

func setKeyValue(key, value string, t *testing.T) {
	t.Helper()
	config := config.NewConfig()

	err := config.SetConfigKeyValue(key, value)
	if err != nil {
		t.Fatalf("Error setting config key: %v", err)
	}
}
