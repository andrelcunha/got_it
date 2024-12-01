package models

import (
	"fmt"
	"got_it/internal/logger"
	"strings"
)

type CommitData struct {
	Tree           string
	Parent         string
	AuthorName     string
	AuthorEmail    string
	AuthorDate     string
	CommitterName  string
	CommitterEmail string
	CommitterDate  string
	Message        string
}

type CommitDataParser struct {
	logger *logger.Logger
}

type CommitKey string

const (
	TREE      CommitKey = "tree"
	PARENT    CommitKey = "parent"
	AUTHOR    CommitKey = "author"
	COMMITTER CommitKey = "committer"
	MESSAGE   CommitKey = "message"
)

var commitKeys = map[CommitKey]string{
	TREE:      "tree",
	PARENT:    "parent",
	AUTHOR:    "author",
	COMMITTER: "committer",
	MESSAGE:   "message",
}

func NewCommitDataParser(logger *logger.Logger) CommitDataParser {
	return CommitDataParser{
		logger: logger,
	}
}

func (cp *CommitDataParser) Parse(commitMetadata string) (CommitData, error) {
	commitData := &CommitData{}
	commitData, error := parseCommitMetadata(cp.logger, commitMetadata, commitData)
	if error != nil {
		cp.logger.Debug("Error parsing commit metadata: %s", error)
		return *commitData, nil
	}
	return *commitData, nil
}

// fetchValueFromKeyInCommitMetadata returns the value of the key in the commit metadata
func parseCommitMetadata(logger *logger.Logger, commitContent string, cd *CommitData) (*CommitData, error) {
	// read the commit metadata, line by line until find the line with 'tree' prefix
	lines := strings.Split(commitContent, "\n")
	fieldsFilled := 0
	flagNextlineIsMessage := false
	var commitMessage string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		logger.Debug("Line: %s", line)

		parts := strings.Split(line, " ")
		if len(parts) >= 2 && fieldsFilled < 4 {
			value := strings.Join(parts[1:], " ")
			switch parts[0] {

			case commitKeys[TREE]:
				cd.Tree = value
				fieldsFilled++

			case commitKeys[PARENT]:
				cd.Parent = value
				fieldsFilled++

			case commitKeys[AUTHOR]:
				name, email, date, err := parseAuthoOrCommiterLine(line)
				if err != nil {
					return nil, err
				}
				cd.AuthorName = name
				cd.AuthorEmail = email
				cd.AuthorDate = date
				fieldsFilled++

			case commitKeys[COMMITTER]:
				name, email, date, err := parseAuthoOrCommiterLine(line)
				if err != nil {
					return nil, err
				}
				cd.CommitterName = name
				cd.CommitterEmail = email
				cd.CommitterDate = date
				fieldsFilled++

			default:
				continue
			}
		}
		// if fieeldsFilled = 4 means the next should be the message
		if line == "" {
			flagNextlineIsMessage = true
			logger.Debug("Next line should be the message")
			continue
		}
		if flagNextlineIsMessage {
			logger.Debug("Getting commit message: %s", line)
			commitMessage += line + "\n"
			continue
		}

	}
	if commitMessage != "" {
		logger.Debug("Commit message: %s", commitMessage)
		cd.Message = strings.TrimSpace(commitMessage)
	}
	return cd, nil
}

// GetKeyCommitMetadata returns the value of the key in the commit metadata
func GetKeyCommitMetadata(logger *logger.Logger, commitContent string, key CommitKey) (string, error) {
	// read the commit metadata, line by line until find the line with 'tree' prefix
	lines := strings.Split(commitContent, "\n")
	flagNextlineIsMessage := false
	var commitMessage string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		logger.Debug("Line: %s", line)
		if line == "" {
			flagNextlineIsMessage = true
			logger.Debug("Next line should be the message")
			continue
		}
		if flagNextlineIsMessage && key == MESSAGE {
			logger.Debug("Getting commit message: %s", line)
			commitMessage += line + "\n"
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) >= 2 {
			logger.Debug("Field: %s", parts[0])
			if parts[0] == commitKeys[key] {
				return strings.Join(parts[1:], " "), nil
			}
		}
	}
	if key == MESSAGE {
		logger.Debug("Commit message: %s", commitMessage)
		return commitMessage, nil
	}
	return "", fmt.Errorf("key %s not found in commit", commitKeys[key])
}

func parseAuthoOrCommiterLine(line string) (string, string, string, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid line: %s", line)
	}
	emailIndexStarts := -1
	emailIndexEnds := -1

	for i, part := range parts {
		if strings.HasPrefix(part, "<") {
			emailIndexStarts = i

		}
		if strings.HasSuffix(part, ">") {
			emailIndexEnds = i
			break
		}
	}

	if emailIndexStarts == -1 || emailIndexEnds == -1 {
		return "", "", "", fmt.Errorf("invalid line: %s", line)
	}
	email := strings.Join(parts[emailIndexStarts:emailIndexEnds+1], " ")
	email = strings.Trim(email, "<>")
	name := strings.Join(parts[1:emailIndexStarts], " ")
	date := strings.Join(parts[emailIndexEnds+1:], " ")

	return name, email, date, nil
}
