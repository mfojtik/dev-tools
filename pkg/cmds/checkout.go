package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mfojtik/dev-tools/pkg/api"
	"github.com/mfojtik/dev-tools/pkg/github"
)

func CheckoutPullRequest(user, pullName string) error {
	client := github.Login()
	// For lazy people, support using the github urls
	if strings.Contains(pullName, "github.com") {
		parts := strings.Split(pullName, "/")
		pullName = parts[len(parts)-1]
	}
	number, err := strconv.ParseInt(pullName, 10, 64)
	if err != nil {
		return err
	}
	pr, _, err := client.PullRequests.Get(api.OriginRepoOwner, api.OriginRepoName, int(number))
	if err != nil {
		return nil
	}
	fmt.Fprintf(os.Stdout, "--> The local branch name for #%d is %q\n", number, *pr.Head.Ref)
	result, err := exec.Command("git", "rev-parse", "--verify", *pr.Head.Ref).Output()
	if err != nil {
		return fmt.Errorf("Unable to find local branch %q: %v", *pr.Head.Ref, err)
	}
	fmt.Fprintf(os.Stdout, "--> Branch %q points to %q\n", *pr.Head.Ref, string(result))
	if _, err := exec.Command("git", "checkout", *pr.Head.Ref).Output(); err != nil {
		return fmt.Errorf("Unable to checkout  branch %q: %v", *pr.Head.Ref, err)
	}
	fmt.Fprintf(os.Stdout, "--> Switched branch to %q\n", *pr.Head.Ref)
	return nil
}
