package commit

import (
	"bufio"
	"fmt"
	"got_it/internal/commands/config"
	"got_it/internal/logger"
	"got_it/internal/models"
	"got_it/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var verbose bool = false

type Commit struct {
	conf       *config.Config
	commitData *models.CommitData
	logger     *logger.Logger
}

func NewCommit(message string) *Commit {
	conf := config.NewConfig()
	l := logger.NewLogger(verbose)
	commitData := &models.CommitData{
		Message: message,
	}

	return &Commit{
		conf:       conf,
		commitData: commitData,
		logger:     l,
	}
}

// Execute is the entry point for the commit command
// It is a shortcut for Commit.NewCommit(message).RunCommit()
func Execute(message string, beVerbose bool) {
	verbose = beVerbose
	co := NewCommit(message)
	co.RunCommit()
}

func (co *Commit) RunCommit() (string, error) {
	err := co.FetchTree()
	if err != nil {
		fmt.Println("Error fetching tree:", err)
		return "", err
	}

	err = co.FetchParent()
	if err != nil {
		co.logger.Debug("Error fetching parent: %s", err)
	}

	err = co.FetchAuthorData(co.commitData)
	if err != nil {
		fmt.Println("Error fetching author data:", err)
		return "", err
	}

	err = co.FetchCommitterData(co.commitData)
	if err != nil {
		fmt.Println("Error fetching committer data:", err)
		return "", err
	}

	commitMetadata := co.FormatCommitMetadata(co.commitData)
	co.logger.Log(commitMetadata)

	// Hash the commit metadata
	commitHash := utils.HashContent(commitMetadata)

	err = co.storeObject(commitHash, commitMetadata)
	if err != nil {
		fmt.Println("Error storing commit object:", err)
		return "", err
	}

	return commitMetadata, co.updateHEAD(commitHash)
}

func (co *Commit) FetchTree() error {
	stagedFiles, err := co.readStagedFiles()
	if err != nil {
		return err
	}
	tree, err := co.generateTreeObject(stagedFiles)
	if err != nil {
		return err
	}
	co.commitData.Tree = tree
	return nil
}

func (co *Commit) FetchParent() error {
	parent, err := co.getParentCommitHash()
	if err != nil {
		return err
	}
	co.commitData.Parent = parent
	return nil
}

func (co *Commit) FetchAuthorData(commitData *models.CommitData) error {
	authorName := co.conf.GetUserName()
	authorEmail := co.conf.GetUserEmail()
	timestamp := time.Now()

	commitData.AuthorName = authorName
	commitData.AuthorEmail = authorEmail
	commitData.AuthorDate = timestamp.Format(time.RFC3339)
	return nil
}

func (co *Commit) FetchCommitterData(commitData *models.CommitData) error {
	committerName := co.conf.GetUserName()
	committerEmail := co.conf.GetUserEmail()
	timestamp := time.Now()

	commitData.CommitterName = committerName
	commitData.CommitterEmail = committerEmail
	commitData.CommitterDate = timestamp.Format(time.RFC3339)
	return nil
}

func (co *Commit) FormatCommitMetadata(commitData *models.CommitData) string {
	var commitStr string

	// Tree
	commitStr += fmt.Sprintf("tree %s\n", commitData.Tree)
	// Parent
	if commitData.Parent != "" {
		commitStr += fmt.Sprintf("parent %s\n", commitData.Parent)
	}
	// Author
	commitStr += fmt.Sprintf("author %s <%s> %s\n", commitData.AuthorName, commitData.AuthorEmail, commitData.AuthorDate)
	// Committer
	commitStr += fmt.Sprintf("committer %s <%s> %s\n", commitData.CommitterName, commitData.CommitterEmail, commitData.CommitterDate)
	// Empty line
	commitStr += "\n"
	// Message
	commitStr += fmt.Sprintf("%s\n", commitData.Message)
	co.logger.Log("Commit metadata:\n\n" + commitStr)
	return commitStr
}

// readStagedFiles opens index file and reads the staged files, returning a map with file names and ther hashes
func (co *Commit) readStagedFiles() (map[string]string, error) {
	stagedFiles := make(map[string]string)

	// Open the index file
	indexFile, err := os.Open(co.conf.GetIndexPath())
	if err != nil {
		return nil, err
	}
	defer indexFile.Close()

	// Read the index file
	scanner := bufio.NewScanner(indexFile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			stagedFiles[parts[0]] = parts[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return stagedFiles, err
	}
	co.logger.Log("Staged files: " + fmt.Sprintf("%v", stagedFiles))

	return stagedFiles, nil
}

// generateTreeObject creates a tree object from the staged files
// It receives a map with file names and their hashes and returns a string with the tree object,
func (co *Commit) generateTreeObject(stagedFiles map[string]string) (string, error) {

	prefix, _ := filepath.Abs(".")
	prefix += separator()
	treeContent := co.generateTreeContent(stagedFiles, prefix)
	treeHash := utils.HashContent(treeContent)
	err := co.storeObject(treeHash, treeContent)

	co.logger.Log("Tree hash: \n\n" + treeHash)
	return treeHash, err
}

// generateTreeContent creates a tree object from the staged files
func (co *Commit) generateTreeContent(stagedFiles map[string]string, prefix string) string {
	var treeContent strings.Builder
	directories := make(map[string]map[string]string)

	// Logging initial staged files map
	co.logger.Log("Initial staged files: %v", stagedFiles)

	for filePath, hash := range stagedFiles {
		relativePath := strings.TrimPrefix(filePath, prefix)
		co.logger.Log("Relative path: %s", relativePath)
		// Get the separator character for the current OS
		parts := strings.SplitN(relativePath, separator(), 2)

		if len(parts) == 1 {
			// It's a file (blob)
			mode, err := co.getFileMode(filePath)
			if err != nil {
				fmt.Printf("Error getting file mode for %s: %v\n", filePath, err)
				continue
			}
			entry := fmt.Sprintf("%s blob %s\t%s\n", mode, hash, relativePath)
			co.logger.Log("File entry: %s\n", entry)
			treeContent.WriteString(entry)
		} else {
			// It's a directory (tree)
			dir := parts[0]
			if directories[dir] == nil {
				directories[dir] = make(map[string]string)
			}
			co.logger.Log("Directory: %s, Remaining Path: %s", dir, parts[1])
			directories[dir][parts[1]] = hash
		}
	}

	for dir, files := range directories {
		co.logger.Log("Processing directory: %s", dir)
		prefix := prefix + dir + separator()
		prefixedFiles := make(map[string]string)
		for file, hash := range files {
			prefixedFiles[prefix+file] = hash
		}
		subTreeContent := co.generateTreeContent(prefixedFiles, prefix)
		subTreeHash := utils.HashContent(subTreeContent)
		co.logger.Log("SubTree Hash: %s for Directory: %s", subTreeHash, dir)
		co.storeObject(subTreeHash, subTreeContent)
		treeContent.WriteString(fmt.Sprintf("040000 tree %s\t%s\n", subTreeHash, dir))
		treeContent.WriteString(subTreeContent)
	}
	return treeContent.String()
}

func (co *Commit) getFileMode(file string) (string, error) {
	var gotMode string
	info, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	mode := info.Mode()
	if info.IsDir() {
		gotMode = "40000"
	} else if mode&0111 != 0 {
		gotMode = "100755"
	} else {
		gotMode = "100644"
	}
	return gotMode, nil
}

// getParentCommitHash returns the hash of the parent commit (the HEAD commit)
func (co *Commit) getParentCommitHash() (string, error) {

	headRef, err := co.readRefFromHEAD()
	if err != nil {
		return "", err
	}

	// Verrify if the file exists
	if _, err := os.Stat(headRef); os.IsNotExist(err) {
		co.logger.Debug("File does not exist: %s", headRef)
		return "", err
	}

	// Read the content of the file pointed to by the HEAD reference
	commitHashBytes, err := os.ReadFile(headRef)
	if err != nil {
		co.logger.Debug("Error reading commit file: %s", err)
		return "", err
	}

	return string(commitHashBytes), nil
}

func (co *Commit) readRefFromHEAD() (string, error) {
	// Get the current commit from the HEAD
	headPath := filepath.Join(co.conf.GotDir, "HEAD")
	headRefBytes, err := os.ReadFile(headPath)
	if err != nil {
		co.logger.Debug("Error reading HEAD file: %s", err)
		return "", err
	}
	headRef := string(headRefBytes)

	//find the prefix "refs: " in  the headRef
	if !strings.HasPrefix(string(headRef), "ref: ") {
		co.logger.Debug("unespected HEAD format: %s", headRef)
		return "", err
	}

	//remove the prefix "ref: "
	headRef = strings.TrimSpace(headRef)
	headRef = headRef[5:]
	headRef = filepath.Join(co.conf.GotDir, headRef)

	return headRef, nil
}

func (co *Commit) storeObject(hash, content string) error {
	objectPath := filepath.Join(co.conf.GotDir, "objects", hash[:2], hash[2:])
	err := os.MkdirAll(filepath.Dir(objectPath), 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return err
	}
	return os.WriteFile(objectPath, []byte(content), 0644)
}

func (co *Commit) generateCommitFeedback(commitMetadata string) string {
	feedback := fmt.Sprintf("Commit details:\n%s\n", commitMetadata)
	return feedback
}

func (co *Commit) updateHEAD(commitHash string) error {
	// GEt the ref to the HEAD file
	headRef, err := co.readRefFromHEAD()
	if err != nil {
		co.logger.Debug("Error reading HEAD file: %s", err)
		return err
	}
	err = os.MkdirAll(filepath.Dir(headRef), 0755)
	f, err := os.Create(headRef)
	f.Close()
	return os.WriteFile(headRef, []byte(commitHash), 0644)
}
