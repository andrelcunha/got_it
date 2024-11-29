package add

import (
	"crypto/rand"
	"fmt"
	"got_it/internal/commands/config"
	"got_it/internal/logger"
	"got_it/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test isSubPath
func TestIsSubPath(t *testing.T) {
	tests := []struct {
		dir  string
		file string
		want bool
	}{
		{"testdata", "testdata/file1.txt", true},
		{"testdata", "testdata/dir1/file1.txt", true},
		{"testdata", "testdata/dir1/dir11/file1.txt", true},
		{".git", ".gitignore", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s is subpath of %s", tt.file, tt.dir), func(t *testing.T) {
			got := isSubPath(tt.dir, tt.file)
			if got != tt.want {
				t.Errorf("isSubPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEssentialFile(t *testing.T) {
	tests := []struct {
		file string
		want bool
	}{
		{".gotignore", true},
		{".gotignore_test", true},
		{"testdata/file1.txt", false},
		{"testdata/dir1/file1.txt", false},
		{"testdata/dir1/dir11/file1.txt", false},
	}
	essentialFiles := []string{
		".gotignore",
		".gotignore_test",
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s is essential file", tt.file), func(t *testing.T) {
			got := isEssentialFile(tt.file, essentialFiles)
			if got != tt.want {
				t.Errorf("isEssentialFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsGotDir checks if the .got directory exists
func TestIsGotDir(t *testing.T) {
	tests := []struct {
		dir  string
		want bool
	}{
		{".got/teste", true},
		{"testdata", false},
	}

	c := config.NewConfig()
	l := logger.NewLogger(false, true)
	a := NewAdd(c, l)
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s is .got directory", tt.dir), func(t *testing.T) {
			gotDir := tt.dir + string(filepath.Separator)
			got := a.isGotDir(gotDir)
			if got != tt.want {
				t.Errorf("isGotDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestUpdateIndexFile updates the index file with the staged files
func TestUpdateIndexFile(t *testing.T) {
	// ARRANGE
	files := []string{
		"testdata/file1.txt",
		"testdata/dir1/file2.txt"}

	t.Setenv("GOT_DEBUG", "true")
	c := config.NewConfig()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	indexFile := c.GetIndexPath()
	err := os.MkdirAll(filepath.Dir(indexFile), 0755)
	if err != nil {
		t.Fatalf("Error creating directory: %v", err)
	}
	f, err := os.Create(indexFile)
	if err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	f.Close()

	defer os.RemoveAll(filepath.Dir(indexFile))
	l := logger.NewLogger(false, true)
	a := NewAdd(c, l)

	for _, file := range files {
		_, hash := createRandomFileGetHash(t, file)
		t.Logf("\nfile: %s\nhash: %s", file, hash)
		addToIndex(indexFile, file, hash)
	}
	//newContent := generateRandomContent(t)
	writeToFile(t, files[0], "Hi")
	newHash := utils.HashContent("Hi")
	t.Logf("\nfile: %s\nnew hash: %s", files[0], newHash)
	// ACT
	err = a.updateHashChangedFileInIndex(files[0], newHash)
	if err != nil {
		t.Fatalf("Error updating index file: %v", err)
	}

	// ASSERT
	// find file1.txt in index file and check if the hash is the same as the new hash
	indexFileContent, err := os.ReadFile(indexFile)
	if err != nil {
		t.Fatalf("Error reading index file: %v", err)
	}
	// read line by line to check if the hash is the same as the new hash
	lines := strings.Split(string(indexFileContent), "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "testdata/file1.txt") {
			// separate the line by spaces
			parts := strings.Split(line, " ")
			// check if the hash is the same as the new hash
			if parts[1] != newHash {
				t.Errorf("Hash is not the same as the new hash")
				t.Logf("Expected: %s", newHash)
				t.Logf("Got: %s", parts[1])
			}
			found = true
			break
		}
	}
	if !found {
		t.Errorf("File not found in index file")
	}
	t.Logf("Index file content: \n%s", string(indexFileContent))
}

// write the content to the file
func writeToFile(t *testing.T, file string, content string) error {
	t.Helper()
	os.MkdirAll(filepath.Dir(file), 0755)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			t.Log("File does not exist and was not created: ", file)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Log("File created successfully")
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func generateRandomContent(t *testing.T) string {
	t.Helper()
	randomBytes := make([]byte, 1024)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return string(randomBytes)
}

func createRandomFileGetHash(t *testing.T, filePath string) (string, string) {
	t.Helper()
	filePath, _ = filepath.Abs(filePath)

	// Write some content to the file
	content := "Hello, world! " // + generateRandomContent(t)
	err := writeToFile(t, filePath, content)
	if err != nil {
		t.Errorf("Error writing to file: %v", err)
	}
	return filePath, utils.HashContent(content)
}
