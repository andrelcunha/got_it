package init

import (
	"got_it/internal/commands/config"
	"os"
	"path/filepath"
	"testing"
)

// Entry point for testing Init
func TestInit(t *testing.T) {
	i := NewInit()

	// Create a temporary directory for the test
	tempDir := t.TempDir()
	t.Log("Temp dir: ", tempDir)
	//change to the temporary directory
	os.Chdir(tempDir)

	i.InitRepo()

	// ASSERT:
	c := i.conf
	testGenerateHEADFile(t, c)

}

// test generateHEADfile
func testGenerateHEADFile(t *testing.T, c *config.Config) {
	// Check if the HEAD file exists
	absGotDir, _ := filepath.Abs(c.GetGotDir())
	if _, err := os.Stat(absGotDir + "/HEAD"); os.IsNotExist(err) {
		t.Errorf("HEAD file does not exist")
	}
}
