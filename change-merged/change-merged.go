package main

import (
	"flag"
	"github.com/vil1/gerra/config"
	"github.com/vil1/gerra/gerrit"
	"github.com/vil1/gerra/jira"
)

var (
	change,
	changeUrl,
	project,
	branch,
	submitter,
	commit string
)

func init() {
	flag.StringVar(&change, "change", "", "")
	flag.StringVar(&changeUrl, "change-url", "", "")
	flag.StringVar(&project, "project", "", "")
	flag.StringVar(&branch, "branch", "", "")
	flag.StringVar(&submitter, "submitter", "", "")
	flag.StringVar(&commit, "commit", "", "")
	flag.Parse()
}

func main() {
	defer func() {
		if err := config.LogFile.Close(); err != nil {
			panic(err)
		}
	}()
	key := gerrit.DefaultClient.GetIssueFromCommit(project, commit)
	if key != "" {
		jira.DefaultClient.Accept(key)
	}
}
