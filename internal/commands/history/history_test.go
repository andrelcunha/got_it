package history

import (
	"got_it/internal/commands/config"
	"got_it/internal/logger"
	"got_it/internal/utils"
	"os"
	"path/filepath"
	"testing"
)

var mockCommitContent string = `
tree 4b825dc642cb6eb9a060e54bf8d69288fbee4904
parent 3b18e56521f7048abf1ab774cfb4f882c7e61fe4
author John Doe <johndoe@example.com> 1623501234 +0200
committer John Doe <johndoe@example.com> 1623501234 +0200

Initial commit with file.txt and README.md
`

func TestGetContentFromHash(t *testing.T) {
	logger := logger.NewLogger(false, false)
	conf := config.NewConfig()
	commitHash := utils.HashContent(mockCommitContent)
	// create a temporary directory for the test
	tmpDir := t.TempDir()
	err := os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}

	// get the first 2 characters of the commit hash
	dirName := commitHash[:2]
	filename := commitHash[2:]
	filePath := filepath.Join(tmpDir, conf.GotDir, "objects", dirName, filename)
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		t.Fatalf("Error creating directory: %v", err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	file.WriteString(mockCommitContent)
	file.Close()
	if err != nil {
		t.Fatalf("Error hashing content: %v", err)
	}

	// ACT
	retrivedContent, err := getContentFromHash(conf, logger, commitHash)

	// ASSERT
	if err != nil {
		t.Fatalf("Error getting content from hash: %v", err)
	}
	if retrivedContent != mockCommitContent {
		t.Fatalf("Expected %s, got %s", mockCommitContent, retrivedContent)
	}

}

// TestReconstructFileContent tests the reconstructFileContent function
func TestReconstructFileContent(t *testing.T) {
	// logger := logger.NewLogger(false, false)
	// conf := config.NewConfig()
	// commitHash := utils.HashContent(mockCommitContent)
	// // create a temporary directory for the test
	// tmpDir := t.TempDir()
}

func TestFindFileInTree(t *testing.T) {
	treeContent :=
		`100644 blob <hash_of_text1.txt> text1.txt
040000 tree <hash_of_subdir1> subdir1
100644 blob <hash_of_subdir1/text1.txt> text1.txt
100644 delta 3b18e56521f7048abf1ab774cfb4f882c7e61fe4 text2.txt
040000 tree <hash_of_subdir1/subdir1> subdir1
100644 blob 4b825dc642cb6eb9a060e54bf8d69288fbee4904 text1.txt`
	fileName := "subdir1/subdir1/text1.txt"
	expectedHash := "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
	expectedType := "blob"
	assertFindFileInTree(t, treeContent, fileName, expectedHash, expectedType)

	fileName = "subdir1/text2.txt"
	expectedHash = "3b18e56521f7048abf1ab774cfb4f882c7e61fe4"
	expectedType = "delta"

}

func assertFindFileInTree(t *testing.T, treeContent string, fileName string, expectedHash string, expectedType string) {
	hash, ttype, err := findFileInTree(treeContent, fileName)
	if err != nil {
		t.Fatalf("Error finding file in tree: %v", err)
	}
	if hash != expectedHash || ttype != expectedType {
		t.Fatalf("Expected %s, got %s", expectedHash, hash)
	}
}
