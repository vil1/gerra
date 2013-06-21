package gerrit

import (
	"testing"
)

func TestGetIssueKey(t *testing.T){
	commit := &Commit{Subject : "[CMS-998] this is a test commit"}
	if commit.GetIssueKey() != "CMS-998" {
		t.Fail()
	}

	commit = &Commit{Subject : "this is a test commit"}
	if commit.GetIssueKey() != "" {
		t.Fail()
	}
}


