package commit

import (
	"bufio"
	"fmt"
	"got_it/internal/commands/config"
	"got_it/internal/models"
	"os"
	"strings"
	"time"
)

type Commit struct {
	conf       *config.Config
	commitData *models.CommitData
}

func NewCommit(message string) *Commit {
	conf := config.NewConfig()
	commitData := &models.CommitData{
		Message: message,
	}

	return &Commit{
		conf:       conf,
		commitData: commitData,
	}
}

// Execute is a shortcut for creating a new commit and running it
func Execute(message string) {
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
		fmt.Println("Error fetching parent:", err)
		return "", err
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

	return commitMetadata, nil
}

func (co *Commit) FetchTree() error {
	// tree := co.conf.GetTree()
	// commitData.Tree = tree
	return nil
}

func (co *Commit) FetchParent() error {
	// parent := co.conf.GetParent()
	// commitData.Parent = parent
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
	commitStr += fmt.Sprintf("parent %s\n", commitData.Parent)
	// Author
	commitStr += fmt.Sprintf("author %s <%s> %s\n", commitData.AuthorName, commitData.AuthorEmail, commitData.AuthorDate)
	// Committer
	commitStr += fmt.Sprintf("committer %s <%s> %s\n", commitData.CommitterName, commitData.CommitterEmail, commitData.CommitterDate)
	// Empty line
	commitStr += "\n"
	// Message
	commitStr += fmt.Sprintf("\n%s\n", commitData.Message)
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

	return stagedFiles, nil
}
