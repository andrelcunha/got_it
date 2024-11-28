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

// Test GetUserAndEmail function
func TestGetUserAndEmail(t *testing.T) {
	// Create a new Commit instance

	shouldBeUser := "testuser"
	shouldBeEmail := "test@example.com"

	// Set Got environment:
	arrangeEnvironment(t, shouldBeUser, shouldBeEmail)
	//generate files and tree content
	addedFiles := generateProceduralFilesAndDirs(t)
	// addd files to repo
	os.Chdir(originalDir)
	add.Execute(addedFiles, true)

	commit := NewCommit("test commit")
	commitMetadata, err := commit.RunCommit()
	if err != nil {
		t.Errorf("Error running commit: %v", err)
	}
	if commitMetadata == "" {
		t.Errorf("Commit metadata is empty")
	}

	// Check for lines with prefixes `tree` and `parent`
	if !strings.Contains(commitMetadata, "tree ") {
		t.Errorf("Commit metadata does not contain a 'tree' line")
	}
	if !strings.Contains(commitMetadata, "parent ") {
		t.Errorf("Commit metadata does not contain a 'parent' line")
	}
	// Check if author and committer are as expected without timestamps
	expectedAuthorPrefix := fmt.Sprintf("author %s <%s>", shouldBeUser, shouldBeEmail)
	expectedCommitterPrefix := fmt.Sprintf("committer %s <%s>", shouldBeUser, shouldBeEmail)
	if !strings.Contains(commitMetadata, expectedAuthorPrefix) {
		t.Errorf("Commit metadata does not contain expected author prefix")
		t.Errorf("Commit metadata: %s", commitMetadata)
	}
	if !strings.Contains(commitMetadata, expectedCommitterPrefix) {
		t.Errorf("Commit metadata does not contain expected committer prefix")
	}
	// Check if the commit message is as expected
	if !strings.Contains(commitMetadata, "test commit") {
		t.Errorf("Commit metadata does not contain the expected commit message")
	}
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
	// Check if the commit hash is stored in the .got/refs/heads/master file
	masterRefFile := filepath.Join(originalDir, ".got", "refs", "heads", "master")
	masterRefContent, err := os.ReadFile(masterRefFile)
	if err != nil {
		t.Fatalf("Error reading .got/refs/heads/master file: %v", err)
	}
	if string(masterRefContent) != commitHash {
		t.Errorf("Commit hash in .got/refs/heads/master does not match the expected commit hash")
	}
}

// TestReadStagedFiles
func TestReadStagedFiles(t *testing.T) {
	// ARRANGE:
	// Create a new Commit instance
	commit := NewCommit("test commit")
	// Set Got environment:
	// addedFiles, _, err := arrangeEnvironment(t, "testuser", "test@example.com")
	// if err != nil {
	// 	t.Fatalf("Error setting up environment: %v", err)
	// }
	arrangeEnvironment(t, "testuser", "test@example.com")
	addedFiles := generateProceduralFilesAndDirs(t)
	// addd files to repo
	os.Chdir(originalDir)
	add.Execute(addedFiles, true)

	// ACT:
	// Read staged files
	stagedFiles, err := commit.readStagedFiles()
	if err != nil {
		t.Fatalf("Error reading staged files: %v", err)
	}

	// ASSERT:
	// Check if the staged files are as expected
	for _, file := range addedFiles {
		if _, ok := stagedFiles[file]; !ok {
			t.Errorf("File %s is not in the staged files", file)
		}
	}
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
	if treeContent != expectedTreeContent {
		t.Errorf("Tree content is not as expected")
		t.Errorf("Expected:\n %s", expectedTreeContent)
		t.Errorf("Got:\n %s", treeContent)
	}
}

//

// HELPER FUNCTIONS

func arrangeEnvironment(t *testing.T, shouldBeUser string, shouldBeEmail string) {
	// Create a temporary directory as Repository
	createTempDir(t)

	//
	initializeRepo(t)
	// create a file called HEAD
	headPath := filepath.Join(originalDir, ".got", "HEAD")
	headFile, err := os.Create(headPath)
	if err != nil {
		t.Fatalf("Error creating HEAD file: %v", err)
	}
	defer headFile.Close()
	// write "ref: refs/heads/master" on HEAD file
	headFileContent := "ref: refs/heads/master"
	_, err = headFile.WriteString(headFileContent)
	if err != nil {
		t.Fatalf("Error writing to HEAD file: %v", err)
	}

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
