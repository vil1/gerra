package config

import (
	"code.google.com/p/goconf/conf"
	"log"
	"os"
	"strings"
)

var (
	configPath,
	User,
	Pwd,
	Host,
	Port,
	JiraBaseUrl string
	LogFile *os.File
)

func init() {
	var errors Errors
	basePath := os.Args[0]
	configPath = basePath[0:strings.LastIndex(basePath, "/")]
	if cfg, err := conf.ReadConfigFile(configPath + "/hooks.conf"); err == nil {
		if User, err = cfg.GetString("jira", "user"); err != nil {
			errors = append(errors, err)
		}
		if Pwd, err = cfg.GetString("jira", "password"); err != nil {
			errors = append(errors, err)
		}
		if JiraBaseUrl, err = cfg.GetString("jira", "baseUrl"); err != nil {
			errors = append(errors, err)
		}
		if Host, err = cfg.GetString("gerrit", "host"); err != nil {
			errors = append(errors, err)
		}
		if Port, err = cfg.GetString("gerrit", "port"); err != nil {
			errors = append(errors, err)
		}
		if len(errors) > 0 {
			panic(errors)
		}
	} else {
		panic(err)
	}
	var err error
	if LogFile, err = os.OpenFile(configPath + "/hooks.log", os.O_APPEND|os.O_WRONLY, 0600); err == nil {
		log.SetOutput(LogFile)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		panic(err)
	}
}

type Errors []error

func (e Errors) Error() string {
	errStrings := make([]string, len(e))
	for i, err := range e {
		errStrings[i] = err.Error()
	}
	return strings.Join(errStrings, "\n")
}
