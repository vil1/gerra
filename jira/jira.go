package jira

import (
	"net/http"
	"fmt"
	"log"
	"bytes"
	"encoding/json"
)

type Issue struct {}

type transition struct {
	identified
	named
}

type Client struct {
	baseUrl, user, password string
	*http.Client
}

type request struct {
	method string
	pathTemplate string
	body interface{}
	parameters []interface{}
}

func NewClient(baseUrl, user, password string)(client *Client){
	return &Client{baseUrl, user, password, http.DefaultClient}
}


func (c *Client)do(req *request)(resp *http.Response, err error){
	url := c.baseUrl + fmt.Sprintf(req.pathTemplate, req.parameters...)
	body := bytes.NewBuffer(make([]byte, 1024))
	json.NewEncoder(body).Encode(req.body)
	request, err := http.NewRequest(req.method, url, body)
	request.SetBasicAuth(c.user, c.password)
	request.Header.Set("Content-Type", "application/json")
	return c.Client.Do(request)
}

type issue struct {
	Id string
	Transitions []transition
	Changelog struct {
		Histories [] history
	}
}

type named struct {
	Name string `json:"name"`
}

type identified struct {
	Id string  `json:"id"`
}

type history struct{
	Id string
	Author named
	Items []item
}

type item struct {
	Field , ToString, FromString string
}

type doTransition struct {
	Transition identified `json:"transition"`
}

type fields struct{
	Assignee named `json:"assignee"`
	Resolution named `json:"resolution"`
}

func (c *Client)Reject(key string){
	issue, err := c.getIssue(key)
	if err != nil {
		log.Fatalln(err)
	}
	if id, err := findTransitionId(issue.Transitions, "Reject"); err == nil {
		postResponse, postError := c.do(&request{
			"POST",
			"/issue/%s/transitions",
			doTransition{
				identified{id},
			},
			[]interface{}{key}})
		if postResponse.StatusCode != http.StatusOK || postError != nil {
			log.Fatalln(postError, postResponse)
		}
		author := findLastAuthor(issue.Changelog.Histories, "Code Review")
		if err := c.setAssignee(issue, author); err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(err)
	}
}

type CommunicationError string

func (c CommunicationError)Error()string {
	return string(c)
}

func (c *Client) getIssue(key string)(*issue, error) {
	response, err := c.do(&request{"GET", "/issue/%s?fields=transitions&expand=transitions,changelog", nil, []interface{}{key}})
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil , CommunicationError(fmt.Sprintf("Unable to fetch issue data : %v", response))
	}
	issue := new(issue)
	err = json.NewDecoder(response.Body).Decode(issue)
	return issue, err
}

func (c *Client)setAssignee(issue *issue, newAssignee string) error {
	response, err := c.do(&request{"PUT", "/issue/%s/assignee", named{newAssignee}, []interface{}{issue.Id}})
	if err != nil {
    	return err
	} else if response.StatusCode != http.StatusOK {
    	return CommunicationError(fmt.Sprintf("Unable to set assignee %v", response))
	}
	return nil
}

type TransitionNotFound struct {
	msg string
}

func (tnfe *TransitionNotFound) Error() string{
	return tnfe.msg
}

func findTransitionId(transitions []transition, name string)(string, error){
	for _, t := range transitions {
		if t.Name == name {
			return t.Id, nil
		}
	}
	return "", &TransitionNotFound{"Cannot " + name}
}

func findLastAuthor(hist []history, status string)(string) {
	for _, h := range hist{
		for _, i := range h.Items {
			if i.Field == "status" &&  i.ToString == status {
				return h.Author.Name
			}
		}
	}
	return ""
}
