package github

import (
	"fmt"
	"os"

	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Login() *gh.Client {
	token := os.Getenv("GITHUB_API_KEY")
	if len(token) == 0 {
		fmt.Fprintf(os.Stderr, "Error: You must set GITHUB_API_KEY or use --api-key=<key>\n")
		os.Exit(1)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return gh.NewClient(tc)
}
