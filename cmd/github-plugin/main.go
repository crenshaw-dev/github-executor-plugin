package main

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/github"

	githubexecutor "github.com/crenshaw-dev/github-executor-plugin/internal"

	"golang.org/x/oauth2"
)

func main() {
	agentToken, err := os.ReadFile("/var/run/argo/token")
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	executor := githubexecutor.NewGitHubExecutor(client, string(agentToken))
	http.HandleFunc("/api/v1/template.execute", githubexecutor.GitHubPlugin(&executor))
	err = http.ListenAndServe(":4356", nil)
	if err != nil {
		panic(err.Error())
	}
}
