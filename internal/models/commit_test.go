package models

import (
	"got_it/internal/logger"
	"testing"
)

var mockCommitContent string = `tree 4b825dc642cb6eb9a060e54bf8d69288fbee4904
parent 3b18e56521f7048abf1ab774cfb4f882c7e61fe4
author John Doe <johndoe@example.com> 1623501234 +0200
committer John Doe <johndoe@example.com> 1623501234 +0200

Initial commit with file.txt and README.md`

// TestParseCommitMetadata tests the parseCommitMetadata function
func TestGetKeyCommitMetadata(t *testing.T) {
	logger := logger.NewLogger(false, false)
	commitContent := mockCommitContent
	keysArray := []CommitKey{TREE, PARENT, AUTHOR, COMMITTER, MESSAGE}
	for _, key := range keysArray {
		value, err := GetKeyCommitMetadata(logger, commitContent, key)
		if err != nil {
			t.Errorf("Error parsing commit metadata: %v", err)
		}
		t.Logf("Value for key %s: %s", key, value)
	}
}

// TestParseCommitMetadata tests the parseCommitMetadata function
func TestParseCommitMetadata(t *testing.T) {
	logger := logger.NewLogger(false, false)
	commitContent := mockCommitContent
	expectedCommitData := CommitData{
		Tree:           "4b825dc642cb6eb9a060e54bf8d69288fbee4904",
		Parent:         "3b18e56521f7048abf1ab774cfb4f882c7e61fe4",
		AuthorName:     "John Doe",
		AuthorEmail:    "johndoe@example.com",
		AuthorDate:     "1623501234 +0200",
		CommitterName:  "John Doe",
		CommitterEmail: "johndoe@example.com",
		CommitterDate:  "1623501234 +0200",
		Message:        "Initial commit with file.txt and README.md",
	}
	commitData, err := parseCommitMetadata(logger, commitContent, &CommitData{})
	if err != nil {
		t.Errorf("Error parsing commit metadata: %v", err)
	}
	//t.Logf("Commit data: %v", commitData)
	if commitData.Tree != expectedCommitData.Tree {
		t.Errorf("Expected tree %s, got %s", expectedCommitData.Tree, commitData.Tree)
	}
	if commitData.Parent != expectedCommitData.Parent {
		t.Errorf("Expected parent %s, got %s", expectedCommitData.Parent, commitData.Parent)
	}
	if commitData.AuthorName != expectedCommitData.AuthorName {
		t.Errorf("Expected author name %s, got %s", expectedCommitData.AuthorName, commitData.AuthorName)
	}
	if commitData.AuthorEmail != expectedCommitData.AuthorEmail {
		t.Errorf("Expected author email %s, got %s", expectedCommitData.AuthorEmail, commitData.AuthorEmail)
	}
	if commitData.AuthorDate != expectedCommitData.AuthorDate {
		t.Errorf("Expected author date %s, got %s", expectedCommitData.AuthorDate, commitData.AuthorDate)
	}
	if commitData.CommitterName != expectedCommitData.CommitterName {
		t.Errorf("Expected committer name %s, got %s", expectedCommitData.CommitterName, commitData.CommitterName)
	}
	if commitData.CommitterEmail != expectedCommitData.CommitterEmail {
		t.Errorf("Expected committer email %s, got %s", expectedCommitData.CommitterEmail, commitData.CommitterEmail)
	}
	if commitData.CommitterDate != expectedCommitData.CommitterDate {
		t.Errorf("Expected committer date %s, got %s", expectedCommitData.CommitterDate, commitData.CommitterDate)
	}
	if commitData.Message != expectedCommitData.Message {
		t.Errorf("Expected message \"%s\", got \"%s\"", expectedCommitData.Message, commitData.Message)
	}
}
