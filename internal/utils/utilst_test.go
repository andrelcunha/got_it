package utils

import (
	"os"
	"testing"
)

func Test_HashContent_vs_HashFile(t *testing.T) {
	// create a temporary directory
	tempDir := t.TempDir()
	// change to the temporary directory
	err := os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}
	// create a temporary file with content
	tempFile, err := os.CreateTemp(".", "testfile")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())
	// write content to the temporary file
	content := "Hello, World!"
	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}
	// get hash using HashContent
	hashContent := HashContent(content)
	if err != nil {
		t.Fatalf("Error hashing content: %v", err)
	}
	// get hash using HashFile
	hashFile, err := HashFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Error hashing file: %v", err)
	}
	// compare the hashes
	if hashContent != hashFile {
		t.Fatalf("Hashes do not match: %s != %s", hashContent, hashFile)
	}
}
