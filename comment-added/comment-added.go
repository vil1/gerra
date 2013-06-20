package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"code.google.com/p/goconf/conf"
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
	commit = flag.String("commit","","")
	verified = flag.Int("VRIF", 0, "")
	validated = flag.Int("CRVW", 0, "")
	client *http.Client
	user, pwd string
)

func init() {
	flag.Parse()
	client = new(http.Client)
	if cfg, err := conf.ReadConfigFile("foo.conf"); err == nil {
		user, err = cfg.GetString("default", "user")
		pwd, err = cfg.GetString("default", "password")
	}
	fmt.Println(user)
	fmt.Println(pwd)
}

func main() {

	var message string
	var err error
	if message, err = acquireMessage(*project, *branch, *commit) ; err != nil {
		return
	}

	key := getIssueKey(message)
	if *validated <= 0 {
		reject(key)
	} else {
		accept()
	}
}

func acquireMessage(proj string, br string, cm string)(msg string, err error) {
	return
}

func getIssueKey(msg string)(key string){
	return
}

type transition struct {
	id string
}

func reject(key string) {
	var request *http.Request
	var err error
	if request, err = http.NewRequest("GET",strings.Join([]string{baseUrl, api, "issues" , key, "transitions"}, "/"), nil ); err != nil {
		request.SetBasicAuth(user, pwd)
	}
	if resp, err := client.Do(request); err == nil {
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		var transitions []transition
		decoder.Decode(transitions)
		fmt.Println(transitions)
	} else {

	}
	return

}

func accept(){

}
