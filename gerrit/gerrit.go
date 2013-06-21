package gerrit

import(
	"os/exec"
	"encoding/json"
	"bytes"
	"bufio"
	"regexp"
	"fmt"
)

var (
	issueRegexp *regexp.Regexp
)

func init(){
	issueRegexp, _ = regexp.Compile("[A-Z]+-[0-9]+")
}

type Commit struct {
	Project string
	Branch string
	Id string
	Subject string
}

func (c * Commit) GetIssueKey() string {
	return string(issueRegexp.Find([]byte(c.Subject)))
}

func GetCommit(project, id string)(commit *Commit){
	cmd := exec.Command("ssh", "-p", "29418", "kasas@hudson2", "gerrit", "query", "--format=JSON", "commit:" + id, "project:" + project, "limit:1")
	buff := bytes.NewBuffer([]byte{})
	cmd.Stdout = buff
	cmd.Run()
	buf := bufio.NewReader(buff)
	if line, err := buf.ReadBytes('\n'); err == nil {
		bufline := bytes.NewBuffer(line)
		commit = new(Commit)
		if err := json.NewDecoder(bufline).Decode(commit) ; err != nil {
			fmt.Printf("json error :%s\n",err.Error())
		}
	}
	return
}

