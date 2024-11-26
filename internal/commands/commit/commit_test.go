package commit

import (
	"fmt"
	"got_it/internal/commands/config"
	"os"
	"testing"
)

// Test GetUserAndEmail function
func TestGetUserAndEmail(t *testing.T) {
	// Create a new Commit instance
	commit := NewCommit("test commit")

	config := config.NewConfig()
	shouldBeUser := "testuser"
	shouldBeEmail := "test@example.com"

	// Set Got environment:
	// Create a temporary directory for testing
	tempdir := t.TempDir()
	err := os.Chdir(tempdir)
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}
	// Initialize the Got_it repository

	err = config.SetConfigKeyValue("user.name", shouldBeUser)
	if err != nil {
		t.Fatalf("Error setting config key: %v", err)
	}
	config.SetConfigKeyValue("user.email", shouldBeEmail)
	if err != nil {
		t.Fatalf("Error setting config key: %v", err)
	}

	shouldBeUser = config.GetUserName()
	fmt.Printf("user: %s\n", shouldBeUser)

	shouldBeEmail = config.GetUserEmail()
	fmt.Printf("email: %s\n", shouldBeEmail)

	commit.SetConfig(config)
	// Call the GetUserAndEmail function
	user, email := commit.GetUserAndEmail()

	// Check if the returned values are correct
	if user != shouldBeUser || email != shouldBeEmail {
		t.Errorf("GetUserAndEmail() = (%s, %s), want (%s, %s)", user, email, shouldBeUser, shouldBeEmail)
	}
}
