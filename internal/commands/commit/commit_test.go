package commit

import (
	"fmt"
	"got_it/internal/commands/add"
	"got_it/internal/commands/config"
	init_ "got_it/internal/commands/init"
	"got_it/internal/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var originalDir string

// TestCommit is the entry point for the commit command
func TestCommit(t *testing.T) {
	// ARRANGE:
	shouldBeUser := "testuser"
	shouldBeEmail := "test@example.com"

	arrangeEnvironment(t, shouldBeUser, shouldBeEmail)
	addedFiles := generateProceduralFilesAndDirs(t)
	os.Chdir(originalDir) //ensure we are in the repo root
	add.Execute(addedFiles, true)

	// ACT:
	commit := NewCommit("test commit")
	defaultBranch := commit.conf.GetDefaultBranch()
	// Read staged files
	stagedFiles, err := commit.readStagedFiles()
	if err != nil {
		t.Fatalf("Error reading staged files: %v", err)
	}

	commitMetadata, err := commit.RunCommit()
	if err != nil || commitMetadata == "" {
		t.Errorf("Failed on running commit: %v", err)
	}

	// ASSERT:
	// Check author and committer are set
	testAuthorAndCommitter(t, shouldBeUser, shouldBeEmail, commitMetadata)
	commitHash := testCommitHash(t, commitMetadata)
	if commitHash == "" {
		t.Fatalf("Commit hash is empty")
	}
	testReadStagedFiles(t, addedFiles, stagedFiles)
	testRefContent(t, defaultBranch, commitHash)
}

// Test GenerateTreeObject, GenereateTreeContent and getFileMode
func TestReadTree(t *testing.T) {
	// ARRANGE:
	// Create a new Commit instance
	commit := NewCommit("test commit")
	verbose = true
	// Set Got environment:
	arrangeEnvironment(t, "testuser", "test@example.com")
	addedFiles, expectedTreeContent := generateFilesAndTreContent(t)
	os.Chdir(originalDir)
	add.Execute([]string{"."}, true)

	// ACT:
	// Read tree
	// get staged files
	stagedFiles, err := commit.readStagedFiles()
	// generate tree content
	separator := string(filepath.Separator)
	prefix, _ := filepath.Abs(".")
	prefix += separator
	treeContent := commit.generateTreeContent(stagedFiles, prefix)

	if err != nil {
		t.Fatalf("Error reading tree: %v", err)
	}
	// ASSERT:
	// Check if the tree is as expected
	if treeContent == "" {
		t.Errorf("Tree is nil")
	}
	// Find file name in the tree(string separated by \n)
	treeLines := strings.Split(treeContent, "\n")

	// Check if the tree contains the added files
	for _, file := range addedFiles {
		found := false
		for _, line := range treeLines {
			if strings.Contains(file, line) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("File %s is not in the tree", file)
		}
	}
	// Check if the tree content is as expected
	// Compare the tree content with the expected tree content, line by line
	for i, line := range strings.Split(treeContent, "\n") {
		for _, expectedLine := range strings.Split(expectedTreeContent, "\n") {
			if line == expectedLine {
				break
			}
			if i == len(treeLines)-1 && line != "" {
				t.Errorf("Tree content is not as expected")
				t.Errorf("Expected:\n %s", expectedTreeContent)
				t.Errorf("Got:\n %s", treeContent)
			}
		}
	}

	// if treeContent != expectedTreeContent {
	// 	t.Errorf("Tree content is not as expected")
	// 	t.Errorf("Expected:\n %s", expectedTreeContent)
	// 	t.Errorf("Got:\n %s", treeContent)
	// }
}

// SUB-TESTS:

// testAuthorAndCommitter tests if the author and committer are set correctly
func testAuthorAndCommitter(t *testing.T, shouldBeUser, shouldBeEmail string, commitMetadata string) error {
	// Check if author and committer are as expected without timestamps
	expectedAuthorPrefix := fmt.Sprintf("author %s <%s>", shouldBeUser, shouldBeEmail)
	expectedCommitterPrefix := fmt.Sprintf("committer %s <%s>", shouldBeUser, shouldBeEmail)
	if !strings.Contains(commitMetadata, expectedAuthorPrefix) {
		t.Errorf("Commit metadata does not contain expected author prefix")
		t.Errorf("Commit metadata: %s", commitMetadata)
		return fmt.Errorf("Commit metadata does not contain expected author prefix")
	}
	if !strings.Contains(commitMetadata, expectedCommitterPrefix) {
		t.Errorf("Commit metadata does not contain expected committer prefix")

	}
	return nil
}

func testCommitHash(t *testing.T, commitMetadata string) string {
	// Create commit hash
	commitHash := utils.HashContent(commitMetadata)
	// check if folder with commit hash exists
	commitFolder := filepath.Join(originalDir, ".got", "objects", commitHash[:2])
	if _, err := os.Stat(commitFolder); os.IsNotExist(err) {
		t.Errorf("Commit folder does not exist")
	}
	// Check if file with commit hash as name exists
	commitFile := filepath.Join(commitFolder, commitHash[2:])
	if _, err := os.Stat(commitFile); os.IsNotExist(err) {
		t.Errorf("Commit file does not exist")
	}
	return commitHash
}

// testReadStagedFiles tests if the staged files are as expected
func testReadStagedFiles(t *testing.T, addedFiles []string, stagedFiles map[string]string) {
	// Check if the staged files are as expected
	for _, file := range addedFiles {
		if _, ok := stagedFiles[file]; !ok {
			t.Errorf("File %s is not in the staged files", file)
		}
	}
}

func testRefContent(t *testing.T, defaultBranch, commitHash string) {
	// Check if the commit hash is stored in the .got/refs/heads/ + defautBranch file
	mainRefFile := filepath.Join(originalDir, ".got", "refs", "heads", defaultBranch)
	mainRefContent, err := os.ReadFile(mainRefFile)
	if err != nil {
		t.Fatalf("Error reading .got/refs/heads/%s file: %v", defaultBranch, err)
	}
	if string(mainRefContent) != commitHash {
		t.Errorf("Commit hash in .got/refs/heads/" + defaultBranch + " does not match the expected commit hash")
	}
}

// HELPER FUNCTIONS

func arrangeEnvironment(t *testing.T, shouldBeUser string, shouldBeEmail string) {
	t.Helper()
	// Create a temporary directory as Repository
	createTempDir(t)

	//
	initializeRepo(t)

	// Initialize the repository
	setUserAndEmail(shouldBeUser, shouldBeEmail, t)

	return
}

func createTempDir(t *testing.T) string {
	tempdir := t.TempDir()
	originalDir = tempdir
	err := os.Chdir(tempdir)
	fmt.Println(tempdir)
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}
	return tempdir
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
func addFilesToRepo(t *testing.T) ([]string, string, error) {
	t.Helper()

	// add the file to the index
	// addedFiles := generateProceduralFilesAndDirs(t)
	addedFiles, treeContent := generateFilesAndTreContent(t)

	add.Execute(addedFiles[:], true)
	return addedFiles[:], treeContent, nil
}

func generateTempFiles(t *testing.T) ([]string, error) {
	t.Helper()
	// create a list to store the files to be added
	var addedFiles [10]string
	// create temporary files
	for i := 0; i < len(addedFiles); i++ {
		tempFile, err := os.CreateTemp(".", "testfile")
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())
		if err != nil {
			return nil, fmt.Errorf("Error creating temporary file: %v", err)
		}
		addedFiles[i] = tempFile.Name()

	}
	return addedFiles[:], nil
}

func generateProceduralFilesAndDirs(t *testing.T) []string {
	t.Helper()
	// create a list to store the files to be added
	template := "testfile00%s.txt"

	var addedFiles []string

	j := 0
	for i := range 9 {
		// if i = 0, 3 or 6 create a directory
		if i%3 == 0 {
			//crate the name of directory dir+j
			dirX := "dir" + strconv.Itoa(j)
			dirX, err := filepath.Abs(dirX)
			if err != nil {
				t.Fatalf("Error creating directory: %v", err)
				return nil
			}
			err = os.Mkdir(dirX, 0755)
			if err != nil {
				t.Fatalf("Error creating directory: %v", err)
				return nil
			}
			t.Cleanup(func() {
				err := os.RemoveAll(dirX)
				if err != nil {
					t.Logf("Warning: Error removing directory: %v", err)
				}
			})
			// change to the new directory
			err = os.Chdir(dirX)
			if err != nil {
				t.Fatalf("Error changing directory: %v", err)
				return nil
			}
			j++
			err = os.Chdir("..")
			if err != nil {
				t.Fatalf("Error changing directory: %v", err)
				return nil
			}
		}
		file, err := os.Create(fmt.Sprintf(template, strconv.Itoa(i)))
		if err != nil {
			t.Fatalf("Error creating temporary file: %v", err)
			return nil
		}
		defer file.Close()
		filePath, err := filepath.Abs(file.Name())
		if err != nil {
			t.Fatalf("Error adding file to list: %v", err)
			return nil
		}
		addedFiles = append(addedFiles, filePath)
	}
	return addedFiles
}

func generateFilesAndTreContent(t *testing.T) ([]string, string) {
	t.Helper()

	// var treeContent strings.Builder
	treeContentList := []string{}
	fileList := []string{
		"file1.txt",
		"file2.txt",
		"subdir/",
		"file3.txt",
	}

	for _, item := range fileList {
		// if the item is a directory, create it and change to it
		if strings.HasSuffix(item, "/") {
			err := os.Mkdir(item, 0755)
			if err != nil {
				t.Fatalf("Error creating directory: %v", err)
			}
			// defer os.RemoveAll(item)
			// chdir to the new directory
			err = os.Chdir(item)
			if err != nil {
				t.Fatalf("Error changing directory: %v", err)
			}
			fileEntry := fmt.Sprintf("040000 tree <hash>\t%s\n", strings.TrimSuffix(item, "/"))
			// treeContent.WriteString(fileEntry)
			treeContentList = append(treeContentList, fileEntry)
			continue
		}
		// if the item is a file, create it and add it to the tree content
		file, err := os.Create(item)
		if err != nil {
			t.Fatalf("Error creating temporary file: %v", err)
		}
		defer file.Close()
		// defer os.Remove(file.Name())
		hash, err := utils.HashFile(file.Name())
		if err != nil {
			t.Fatalf("Error hashing file: %v", err)
		}
		fileEntry := fmt.Sprintf("100644 blob %s\t%s\n", hash, file.Name())
		// treeContent.WriteString(fileEntry)
		treeContentList = append(treeContentList, fileEntry)
	}

	// change back to the original directory
	_ = os.Chdir(originalDir)
	flagDirContent := false
	var dirContent strings.Builder
	for _, item := range treeContentList {
		if strings.HasPrefix(item, "040000") {
			flagDirContent = true
			continue
		}
		if flagDirContent {
			dirContent.WriteString(item)
		}
	}
	hash := utils.HashContent(dirContent.String())

	for i, item := range treeContentList {
		if strings.HasPrefix(item, "040000") {
			// replace "<hash>" with the new hash
			treeContentList[i] = strings.Replace(item, "<hash>", hash, 1)
			break
		}
	}

	var treeContent strings.Builder
	for _, item := range treeContentList {
		treeContent.WriteString(item)
	}
	var addedFiles []string
	prefix := originalDir
	for _, item := range fileList {
		if strings.HasSuffix(item, "/") {
			prefix = filepath.Join(prefix, item)
			continue
		}
		filePath := filepath.Join(originalDir, item)
		addedFiles = append(addedFiles, filePath)
	}

	return addedFiles, treeContent.String()

}
