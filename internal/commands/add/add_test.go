package add

import (
	"fmt"
	"got_it/internal/commands/config"
	"os"
	"path/filepath"
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
	a := NewAdd(false, c)
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
