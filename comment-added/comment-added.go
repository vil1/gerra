package main

import (
	"flag"
	"github.com/vil1/gerra/config"
	"github.com/vil1/gerra/gerrit"
	"github.com/vil1/gerra/jira"
)

const (
	baseUrl = "http://jira.fullsix.com"
	api     = "rest/api/latest"
)

var (
	changeId,
	changeUrl,
	comment,
	project,
	author,
	branch,
	commitId string
	verified,
	validated int
)

func init() {
	flag.StringVar(&changeId, "change", "", "")
	flag.StringVar(&changeUrl, "change-url", "", "")
	flag.StringVar(&comment, "comment", "", "")
	flag.StringVar(&project, "project", "", "")
	flag.StringVar(&author, "author", "", "")
	flag.StringVar(&branch, "branch", "", "")
	flag.StringVar(&commitId, "commit", "", "")
	flag.IntVar(&verified, "VRIF", 0, "")
	flag.IntVar(&validated, "CRVW", 0, "")
	flag.Parse()
}

func main() {
	defer func() {
		if err := config.LogFile.Close(); err != nil {
			panic(err)
		}
	}()
	if validated < 0 || verified < 0 {
		if key := gerrit.DefaultClient.GetIssueFromCommit(project, commitId); key != "" {
			jira.DefaultClient.Reject(key)
		}
	}
}
