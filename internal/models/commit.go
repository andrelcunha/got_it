package models

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
