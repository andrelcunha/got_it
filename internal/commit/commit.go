package commit

import (
	"got_it/internal/config"
	"got_it/internal/models"
)

type Commit struct {
	conf     *config.Config
	message  string
	userData *models.User
}

func NewCommit(message string) *Commit {
	conf := config.NewConfig()
	userData := &models.User{
		User:  conf.GetUserName(),
		Email: conf.GetUserEmail(),
	}
	userData.User = conf.GetUserName()
	return &Commit{
		message:  message,
		conf:     conf,
		userData: userData,
	}
}

func (co *Commit) GetMessage() string {
	return co.message
}

func (co *Commit) GetUserAndEmail() (string, string) {
	return co.userData.User, co.userData.Email
}
