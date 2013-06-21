package main

import (
	"flag"
	"github.com/vil1/gerra/gerrit"
	"github.com/vil1/gerra/jira"
	"github.com/vil1/gerra/config"
	"log"
	"os"
)

const (
	baseUrl = "http://jira.fullsix.com"
	api = "rest/api/latest"
)

var(
	changeId = flag.String("change", "", "")
	changeUrl = flag.String("change-url", "", "")
	comment = flag.String("comment", "", "")
	project = flag.String("project", "", "")
	author = flag.String("author", "", "")
	branch = flag.String("branch", "", "")
	commitId = flag.String("commit","","")
	verified = flag.Int("VRIF", 0, "")
	validated = flag.Int("CRVW", 0, "")
	jiraClient *jira.Client
	gerritClient * gerrit.Client
	logFile *os.File
)

func init() {
	flag.Parse()
	jiraClient = jira.NewClient("http://jira.fullsix.com/rest/api/2", config.User, config.Pwd)
	gerritClient = gerrit.NewClient(config.Host, config.Port)
 	var err error
	if logFile, err = os.OpenFile("hooks.log", os.O_APPEND | os.O_WRONLY, 0600); err == nil {
		log.SetOutput(logFile)
	} else {
		panic(err)
	}
}

func main() {
	defer func(){
		if err := logFile.Close(); err != nil {
			panic(err)
		}
	}()
	commit := gerritClient.GetCommit(*project, *commitId)
	key := commit.GetIssueKey()
	if *validated < 0 || *verified < 0{
		jiraClient.Reject(key)
	}
}
