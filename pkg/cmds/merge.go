package cmds

import (
	"fmt"
	"os"

	gh "github.com/google/go-github/github"
	"github.com/mfojtik/dev-tools/pkg/api"
	"github.com/mfojtik/dev-tools/pkg/github"
)

var mergeComment = "[merge]"

func AddMergeComment(pullId int) error {
	client := github.Login()
	c := &gh.IssueComment{Body: &mergeComment}
	comment, _, err := client.Issues.CreateComment(api.OriginRepoOwner, api.OriginRepoName, pullId, c)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Pull request #%d tagged for merge %q\n", pullId, *comment.HTMLURL)
	return nil
}
