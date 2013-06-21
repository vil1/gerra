package config

import (
	"code.google.com/p/goconf/conf"
	"strings"
)

var (
	User,
	Pwd,
	Host,
	Port string
)

func init(){
	var errors Errors
	if cfg, err := conf.ReadConfigFile("gerra.conf"); err == nil {
		if User, err = cfg.GetString("jira", "user") ; err != nil {
			errors = append(errors, err)
		}
		if Pwd, err = cfg.GetString("jira", "password"); err != nil {
			errors = append(errors, err)
		}
		if Host, err = cfg.GetString("gerrit", "user") ; err != nil {
			errors = append(errors, err)
		}
		if Port, err = cfg.GetString("gerrit", "password"); err != nil {
			errors = append(errors, err)
		}
		if len(errors) > 0 {
			panic(errors)
		}
	} else {
		panic(err)
	}
}


type Errors []error

func (e Errors)Error()string {
	errStrings := make([]string, len(e))
	for i, err := range e {
		errStrings[i] = err.Error()
	}
	return strings.Join(errStrings, "\n")
}
