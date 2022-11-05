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
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	executor := githubexecutor.NewGitHubExecutor(client)
	http.HandleFunc("/api/v1/template.execute", githubexecutor.GitHubPlugin(&executor))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err.Error())
	}
}
