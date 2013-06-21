package main

import (
	"flag"
	"fmt"
	"code.google.com/p/goconf/conf"
	"github.com/vil1/gerra/gerrit"
	"github.com/vil1/gerra/jira"
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
	client *jira.Client
	user, pwd string
	logFile *os.File
)

func init() {
	flag.Parse()
	if cfg, err := conf.ReadConfigFile("foo.conf"); err == nil {
		if user, err = cfg.GetString("jira", "user") ; err != nil {
			panic(err)
		}
		if pwd, err = cfg.GetString("jira", "password"); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	client = jira.NewClient("http://jira.fullsix.com/rest/api/2", user, pwd)
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
	commit := gerrit.GetCommit(*project, *commitId)
	fmt.Printf("commit : %v", commit)
	key := commit.GetIssueKey()
	fmt.Printf("key : %s\n",key)
	//if *validated < 0 || *verified < 0{
		client.Reject(key)
	//}
}
