package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	gh "github.com/google/go-github/github"
	"github.com/mfojtik/dev-tools/pkg/api"
	"github.com/mfojtik/dev-tools/pkg/github"
)

func GetPullRequest(user string, quiet bool) error {
	client := github.Login()
	result, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return err
	}
	if len(user) == 0 {
		return fmt.Errorf("no user specified, use --user option")
	}
	branchName := strings.TrimSpace(string(result))
	opts := &gh.PullRequestListOptions{Head: user + ":" + branchName}
	if !quiet {
		fmt.Fprintf(os.Stdout, "Searching pull request based on %q ...\n", user+":"+branchName)
	}
	prs, _, err := client.PullRequests.List(api.OriginRepoOwner, api.OriginRepoName, opts)
	if err != nil {
		if quiet {
			return nil
		}
		return err
	}
	if len(prs) == 0 && !quiet {
		return fmt.Errorf("no pull request for %q", user+":"+branchName)
	}
	for _, p := range prs {
		if quiet {
			fmt.Fprintf(os.Stdout, "%d\n", *p.Number)
			continue
		}
		fmt.Fprintf(os.Stdout, "[%s] #%d: %q (%s)\n", strings.ToUpper(*p.State), *p.Number, *p.Title, *p.HTMLURL)
	}
	return nil
}
