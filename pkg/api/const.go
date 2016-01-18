package api

import (
	"os/exec"
	"strings"
)

var (
	OriginRepoOwner    = "mfojtik"
	OriginRepoName     = "dev-tools"
	UpstreamRemoteName = "upstream"
	OriginBuilders     = []string{
		"openshift/origin-docker-builder",
		"openshift/origin-sti-builder",
	}
)

func init() {
	remoteUrl, err := exec.Command("git", "remote", "get-url", "--push", UpstreamRemoteName).Output()
	if err == nil && len(remoteUrl) > 0 {
		if parts := strings.Split(strings.TrimSpace(string(remoteUrl)), "/"); len(parts) > 2 {
			OriginRepoName = parts[len(parts)-1]
			OriginRepoOwner = parts[len(parts)-2]
		}
	}
}
