package gerrit

import (
	"bufio"
	"encoding/json"
	"github.com/vil1/gerra/config"
	"log"
	"os/exec"
	"regexp"
)

var (
	issueRegexp   *regexp.Regexp
	DefaultClient *Client
)

func init() {
	issueRegexp, _ = regexp.Compile("[A-Z]+-[0-9]+")
	DefaultClient = &Client{config.Host, config.Port}

}

type Commit struct {
	Project string
	Branch  string
	Id      string
	Subject string
}

func (c *Commit) getIssueKey() string {
	return string(issueRegexp.Find([]byte(c.Subject)))
}

type Client struct {
	host, port string
}

func NewClient(host, port string) *Client {
	return &Client{host, port}
}

func (c *Client) GetIssueFromCommit(project, commit string) (key string) {
	cmd := exec.Command("ssh", "-p", c.port, c.host, "gerrit", "query", "--format=JSON", "commit:"+commit, "project:"+project, "limit:1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	if err = cmd.Start(); err != nil {
		log.Fatalln(err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() && key == "" {
		commit := new(Commit)
		if err := json.Unmarshal(scanner.Bytes(), commit); err == nil {
			key = commit.getIssueKey()
		}
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalln(err)
	}
	return
}
