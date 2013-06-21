package jira

import (
	"net/http"
	"fmt"
	"io"
	"log"
	"bytes"
	"encoding/json"
	"io/ioutil"
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
	body io.Reader
	parameters []interface{}
}

func NewClient(baseUrl, user, password string)(client *Client){
	return &Client{baseUrl, user, password, http.DefaultClient}
}


func (c *Client)do(req *request)(resp *http.Response, err error){
	url := c.baseUrl + fmt.Sprintf(req.pathTemplate, req.parameters...)
	request, err := http.NewRequest(req.method, url, req.body)
	request.SetBasicAuth(c.user, c.password)
	request.Header.Set("Content-Type", "application/json")
	return c.Client.Do(request)
}

type issue struct {
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
	//Fields fields `json:"-"`
	Transition identified `json:"transition"`
}

type fields struct{
	Assignee named `json:"assignee"`
	Resolution named `json:"resolution"`
}

func (c *Client)Reject(key string){
	response, err := c.do(&request{"GET", "/issue/%s?fields=transitions&expand=transitions,changelog", nil, []interface{}{key}})
	if err == nil {
		if response.StatusCode == http.StatusOK {
			var issue  = new(issue)
			if err = json.NewDecoder(response.Body).Decode(issue); err == nil {
				fmt.Println()
				fmt.Println(issue)
		   		if id, err := findTransitionId(issue.Transitions, "Reject"); err == nil {
					buff := bytes.NewBuffer([]byte{})
					json.NewEncoder(buff).Encode(doTransition{
						identified{id},
					})
					postResponse, postError := c.do(&request{"POST", "/issue/%s/transitions", buff, []interface{}{key}})
					if postResponse.StatusCode != http.StatusOK || postError != nil {
						fmt.Println("ERROR on POST :")
						fmt.Println(postError)
						body , _ :=  ioutil.ReadAll(postResponse.Body)
						fmt.Println(string(body))
						fmt.Println("================")
					}
					author := findLastAuthor(issue.Changelog.Histories, "Code Review")
					buff = bytes.NewBuffer([]byte{})
					json.NewEncoder(buff).Encode(named{author})
					postResponse, postError = c.do(&request{"PUT", "/issue/%s/assignee", buff, []interface{}{key}})
					if postResponse.StatusCode != http.StatusOK || postError != nil {
						fmt.Println("ERROR on POST :")
						fmt.Println(postError)
						body , _ :=  ioutil.ReadAll(postResponse.Body)
						fmt.Println(string(body))
						fmt.Println("================")
					}
				} else {
					log.Fatalln(err)
				}
			} else {
			   log.Fatalln(err)
			}

		} else {
			log.Fatalf("Unexpected response from JIRA : %s\n", response.Status)
		}
	} else {
		log.Fatalln(err)
	}
}

func (c *Client)Accept(key string){
	response, err := c.do(&request{"GET", "/issue/%s/transitions?expand=transitions", nil, []interface{}{key}})
	if err == nil {
		if response.StatusCode == http.StatusOK {
			var issue  = new(issue)
			if err = json.NewDecoder(response.Body).Decode(issue); err == nil {
				fmt.Printf("Found issue %v\n", issue)
		   		if id, err := findTransitionId(issue.Transitions, "Submit to QA"); err == nil {
					buff := bytes.NewBuffer([]byte{})
					json.NewEncoder(buff).Encode(doTransition{
						identified{id},
					})
					if resp, er := c.do(&request{"POST", "/issue/%s/transitions", buff, []interface{}{key}}); er != nil || resp.StatusCode != http.StatusOK {
						log.Fatalf("error [%v] Jira responded %s to the request\n",er, resp.StatusCode)
					} else {
						fmt.Printf("JIRA responded %v\n", resp)
					}
				} else {
					log.Fatal(err)
				}
			} else {
				fmt.Println("error while fetching issue data", err)
			   	log.Fatal(err)
			}
		}
	} else {
		log.Fatal(err)
	}
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
