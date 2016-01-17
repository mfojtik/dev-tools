# OpenShift development tools

This repository provides few useful tools that developers can use when developing
[OpenShift](https://github.com/openshift/origin).

## Installation

If you have Go installed, you can install `otp` command using:

```console
$ git clone https://github.com/mfojtik/dev-tools
$ cd dev-tools && make install
```

## Configuration

The `otp` tool is smart enough to quess a lot of things based on your current directory
or the system environment. However, you must tell this utility what Github API token
it should use to authorize requests to Github.

To generate new Github token, navigate to [Settings->Access Tokens](https://github.com/settings/tokens)
page and click "Generate new token". Then copy the generated token into an environment variable
you can append to your `~/.bash_profile` (or any other file with more restrictive permissions):

```console
$ echo "GITHUB_API_TOKEN=<token here>" >> ~/.bashrc
$ source ~/.bashrc
```

## Usage

### get

This command will attempt to guess the pull requests you have opened for the current branch.
The pull requests must be opened against the `upstream` remote (given that you are working on a branch
in your fork). Sample usage:

```console
$ cd dev-tools && git checkout test
$ otp get
Searching pull request based on "mfojtik:test" ...
[OPEN] #1: "Test" (https://github.com/mfojtik/dev-tools/pull/1
```

### test

In OpenShift we use special comments with tags to trigger testing in our CI (Jenkins).
Depending on what test you want to trigger, you can set `--extended`, `--only-extended`
or just simple test. Sample usage:

```console
$ otp test $(otp get -n)
Pull request #1 tagged for test "https://github.com/mfojtik/dev-tools/pull/1#issuecomment-172376669"

$ otp test $(otp get -n) --only-extended --focus "build"
Pull request #1 tagged for test "https://github.com/mfojtik/dev-tools/pull/1#issuecomment-172376669"

$ otp test $(otp get -n) --extended --group networking"
Pull request #1 tagged for test "https://github.com/mfojtik/dev-tools/pull/1#issuecomment-172376669"
```

### merge

We don't merge pull requests in OpenShift directly using Github. Instead we use
our CI (Jenkins) to merge them using a merge queue. Sample usage:

```console
$ otp merge $(otp get -n)
Pull request #1 tagged for merge "https://github.com/mfojtik/dev-tools/pull/1#issuecomment-172376669"
```
