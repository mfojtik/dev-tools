package cmds

import (
	"fmt"

	"github.com/mfojtik/dev-tools/pkg/api"
	"github.com/mfojtik/dev-tools/pkg/github"
)

func AddMergeComment(pullId int) error {
	gh := github.Login()
	pr, _, err := gh.PullRequests.Get(api.OriginRepoOwner, api.OriginRepoName, pullId)
	if err != nil {
		return err
	}
	fmt.Printf("pr=%+v\n", pr)
	return nil
}
