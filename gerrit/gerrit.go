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

// Commit is the internal representation of gerrit's commit informations
type Commit struct {
	Project string // the project name
	Branch  string // the branch this commit is included in
	Id      string // the commit's SHA-1 id
	Subject string // the commit message
}

func (c *Commit) getIssueKey() string {
	return string(issueRegexp.Find([]byte(c.Subject)))
}

type Client struct {
	host, port string
}

// NewClient create a new instance of Client for the given host and port 
// and returns a pointer to it.
func NewClient(host, port string) *Client {
	return &Client{host, port}
}

// GetIssueFromCommit obtains the informations about the commit identified by sha1 in the given project.
// It then extracts the issue key from the commit's message.
// The extracted key is the first substring matching [A-Z]+-[0-9]+ found
// in the commit message. If no key is found, a nil string is returned.
// TODO : properly handle gerrit errors.
func (c *Client) GetIssueFromCommit(project, sha1 string) (key string) {
	cmd := exec.Command("ssh", "-p", c.port, c.host, "gerrit", "query", "--format=JSON", "commit:"+sha1, "project:"+project, "limit:1")
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
