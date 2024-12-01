package history

import (
	"fmt"
	"got_it/internal/commands/config"
	"got_it/internal/logger"
	"got_it/internal/models"
	"log"
)

type History struct {
	conf   *config.Config
	logger *logger.Logger
}

func NewHistory(conf *config.Config, logger *logger.Logger) *History {
	return &History{
		conf:   conf,
		logger: logger,
	}
}

func Execute() {
	conf := config.NewConfig()
	logger := logger.NewLogger(false, false)
	hi := NewHistory(conf, logger)
	hi.traverseCommitHistory()
}

func (hi *History) traverseCommitHistory() {
	parser := models.NewCommitDataParser(hi.logger)
	commitHash, branch, err := GetFirstCommitHash(hi.conf, hi.logger)
	isHead := true
	if err != nil {
		log.Fatalf("Error reading HEAD file: %v", err)
	}
	for commitHash != "" {
		if isHead {
			isHead = false
			fmt.Printf("Commit: %s (HEAD -> %s) \n", commitHash, branch)
		} else {
			fmt.Printf("Commit: %s\n", commitHash)
		}

		commitContent, err := getContentFromHash(hi.conf, hi.logger, commitHash)
		if err != nil {
			log.Fatalf("Error reading commit object: %v", err)
		}
		commitMetadata, err := parser.Parse(commitContent)

		parentHash := commitMetadata.Parent
		commitHash = parentHash
		fmt.Printf("Author: %s <%s>\n", commitMetadata.AuthorName, commitMetadata.AuthorEmail)
		fmt.Printf("Date: %s\n", commitMetadata.CommitterDate)
		fmt.Printf("\n" + `    `)
		fmt.Printf("%s\n", commitMetadata.Message)
	}
}
