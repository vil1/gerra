package gerrit

import (
	"bufio"
	"encoding/json"
	//"github.com/vil1/gerra/config"
	"fmt"
	"io"
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
	//DefaultClient = &Client{config.Host, config.Port}

}

type Request struct {
	args map[string]interface{}
}

func NewRequest() *Request {
	return &Request{}
}

// Change is the internal representation of gerrit's commit informations
type Change struct {
	Project     string // the project name
	Branch      string // the branch this commit is included in
	Id          string // the commit's SHA-1 id
	Subject     string // the commit message
	CreatedOn   json.Number
	LastUpdated json.Number
	Status      string
	Number      string
	Owner       *User
	Open        bool
	Sortkey     string
	Url         string
}

type User struct {
	Name, Email, Username string
}

func (c *Change) GetIssueKey() string {
	return string(issueRegexp.Find([]byte(c.Subject)))
}

type Client struct {
	host, port string
}

func (c *Client) prepareCommand(req *Request) *exec.Cmd {
	arguments := []string{"-p", c.port, c.host, "gerrit", "query", "--format=JSON"}
	for key, value := range req.args {
		arguments = append(arguments, fmt.Sprintf("%s:%v", key, value))
	}
	return exec.Command("ssh", arguments...)

}

// NewClient create a new instance of Client for the given host and port
// and returns a pointer to it.
func NewClient(host, port string) *Client {
	return &Client{host, port}
}

// GetCommit obtains the informations about the commit identified by sha1 in the given project.
// It then extracts the issue key from the commit's message.
// The extracted key is the first substring matching [A-Z]+-[0-9]+ found
// in the commit message. If no key is found, a nil string is returned.
// TODO : properly handle gerrit errors.
func (c *Client) GetIssueFromCommit(project, sha1 string) (key string) {
	if commit, err := c.GetCommit(project, sha1); err == nil {
		key = commit.GetIssueKey()
	}
	return
}

// GetCommit obtains from gerrit the commit identified by sha1 in the given project.
func (c *Client) GetCommit(project, sha1 string) (commit *Change, err error) {
	cmd := exec.Command("ssh", "-p", c.port, c.host, "gerrit", "query", "--format=JSON", "commit:"+sha1, "project:"+project, "limit:1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	if err = cmd.Start(); err != nil {
		log.Fatalln(err)
	}
	commit = new(Change)
	scanner := bufio.NewScanner(stdout)
	scanner.Scan()
	if err = json.Unmarshal(scanner.Bytes(), commit); err != nil {
		return
	}
	err = cmd.Wait()
	return
}

func ReadCommits(in io.Reader) (commits []Change, summary map[string]interface{}) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		commit := Change{}
		data := scanner.Bytes()
		if err := json.Unmarshal(data, &commit); err == nil && commit != (Change{}) {
			commits = append(commits, commit)
		} else {
			summary = make(map[string]interface{})
			json.Unmarshal(data, &summary)
		}
	}
	return
}
