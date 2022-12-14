package plugin

import (
	"context"

	"github.com/google/go-github/github"
)

type GitHubClient struct {
	Issues GitHubIssuesClient
	Checks GitHubChecksClient
}

type GitHubIssuesClient interface {
	CreateComment(ctx context.Context, owner, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
	Create(ctx context.Context, owner, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

type GitHubChecksClient interface {
	CreateCheckRun(ctx context.Context, owner, repo string, opt github.CreateCheckRunOptions) (*github.CheckRun, *github.Response, error)
}
