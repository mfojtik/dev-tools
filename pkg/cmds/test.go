package cmds

import (
	"fmt"
	"os"

	gh "github.com/google/go-github/github"
	"github.com/mfojtik/dev-tools/pkg/api"
	"github.com/mfojtik/dev-tools/pkg/github"
)

var (
	testComment         = "[test]"
	onlyExtendedComment = "[testonlyextended]"
)

func AddTestComment(pullId int, extended, onlyExtended bool, focus, group string) error {
	client := github.Login()
	var commentString = ""
	if extended {
		if len(group) > 0 {
			commentString = testComment + "[extended:" + group + "]"
		} else {
			commentString = testComment + "[extended]"
		}
	}
	if onlyExtended {
		if len(focus) > 0 {
			group += "(" + focus + ")"
		}
		if len(group) > 0 {
			commentString = onlyExtendedComment + "[extended:" + group + "]"
		}
	}
	if len(commentString) == 0 {
		commentString = testComment
	}
	c := &gh.IssueComment{Body: &commentString}
	comment, _, err := client.Issues.CreateComment(api.OriginRepoOwner, api.OriginRepoName, pullId, c)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Test comment added to %q\n", *comment.HTMLURL)
	return nil
}
