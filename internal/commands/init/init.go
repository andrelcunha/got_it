package init

import (
	"fmt"
	"got_it/internal/commands/config"
	"got_it/internal/logger"
	"os"
	"path/filepath"
)

type Init struct {
	conf   *config.Config
	logger *logger.Logger
}

func NewInit() *Init {
	conf := config.NewConfig()
	logger := logger.NewLogger(false)
	return &Init{
		conf:   conf,
		logger: logger,
	}
}

func (i *Init) IsInitialized() bool {
	if _, err := os.Stat(i.conf.GetGotDir()); os.IsNotExist(err) {
		fmt.Println("Not a Got_it repository. Run 'got init' first.")
		return false
	}
	return true
}

func (i *Init) InitRepo() {
	gotDir := i.conf.GetGotDir()
	// get absolute path of gotDir
	gotDir, err := filepath.Abs(gotDir)

	// Check if the .got directory already exists
	if _, err := os.Stat(gotDir); !os.IsNotExist(err) {
		fmt.Println("Repository already initialized.")
		return
	}

	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	// Create the .got directory
	if err := os.Mkdir(gotDir, 0755); err != nil {
		fmt.Println("Error creating .got directory:", err)
		return
	}
	// Create the .got/refs/heads directory
	if err := os.MkdirAll(gotDir+"/refs/heads", 0755); err != nil {
		fmt.Println("Error creating .got/refs/heads directory:", err)
		return
	}
	// Generate HEAD file
	err = i.generateHEADfile(gotDir)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Create the .got/objects directory
	if err := os.Mkdir(gotDir+"/objects", 0755); err != nil {
		fmt.Println("Error creating .got/objects directory:", err)
		return
	}

	fmt.Printf("Initialized empty Git repository in %s\n", gotDir)
}

// Generate HEAD file
func (i *Init) generateHEADfile(gotDir string) error {
	// create a file called HEAD
	headPath := filepath.Join(gotDir, "HEAD")

	i.logger.Debug("HEAD file path: %s", headPath)

	headFile, err := os.Create(headPath)
	if err != nil {
		return fmt.Errorf("Error creating HEAD file: %v\n", err)
		// LOG ERROR
	}
	defer headFile.Close()
	defaultBranch := i.conf.GetDefaultBranch()
	refsPath := filepath.Join("refs", "heads", defaultBranch)
	headFileContent := "ref: " + refsPath
	_, err = headFile.WriteString(headFileContent)
	if err != nil {
		return fmt.Errorf("Error writing to HEAD file: %v\n", err)
	}
	return nil
}
