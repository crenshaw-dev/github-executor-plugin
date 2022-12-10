package plugin

import (
	"context"

	"github.com/google/go-github/github"
)

type GitHubClient struct {
	Issues GitHubIssuesClient
}

type GitHubIssuesClient interface {
	CreateComment(ctx context.Context, owner, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
	Create(ctx context.Context, owner, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}
