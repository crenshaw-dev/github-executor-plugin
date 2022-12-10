package plugin

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/crenshaw-dev/github-executor-plugin/internal/mocks"
)

func Test_validateIssueAction(t *testing.T) {
	t.Run("fails on no valid action", func(t *testing.T) {
		err := validateIssueAction(&IssueActionSpec{})
		assert.Error(t, err)
	})
	t.Run("fails on duplicate actions", func(t *testing.T) {
		err := validateIssueAction(&IssueActionSpec{
			Comment: &CommentAction{},
			Create:  &IssueCreateAction{},
		})
		assert.Error(t, err)
	})
}

func Test_runIssueAction(t *testing.T) {
	issuesClient := mocks.NewGitHubIssuesClient(t)
	var r *github.IssueRequest
	issuesClient.Mock.On("Create", mock.Anything, "test", "test", r).Return(&github.Issue{}, nil, nil)
	client := &GitHubClient{
		Issues: issuesClient,
	}
	e := NewGitHubExecutor(client, "")
	_, expectedStatusCode, err := e.runIssueAction(context.Background(), &IssueActionSpec{
		Create: &IssueCreateAction{
			Owner: "test",
			Repo:  "test",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedStatusCode, 201)
}

func Test_validateIssueCreateCommentAction(t *testing.T) {
	t.Run("fails on empty comment body", func(t *testing.T) {
		_, _, _, _, err := validateIssueCreateCommentAction(&CommentAction{})
		assert.Error(t, err)
	})
	t.Run("fails on empty owner", func(t *testing.T) {
		_, _, _, _, err := validateIssueCreateCommentAction(&CommentAction{
			Body: "test",
		})
		assert.Error(t, err)
	})
	t.Run("fails on empty repo", func(t *testing.T) {
		_, _, _, _, err := validateIssueCreateCommentAction(&CommentAction{
			Body:  "test",
			Owner: "test",
		})
		assert.Error(t, err)
	})
	t.Run("fails on empty number", func(t *testing.T) {
		_, _, _, _, err := validateIssueCreateCommentAction(&CommentAction{
			Body:  "test",
			Owner: "test",
			Repo:  "test",
		})
		assert.Error(t, err)
	})
	t.Run("fails on invalid number", func(t *testing.T) {
		_, _, _, _, err := validateIssueCreateCommentAction(&CommentAction{
			Body:   "test",
			Owner:  "test",
			Repo:   "test",
			Number: "test",
		})
		assert.Error(t, err)
	})
	t.Run("fails on negative number", func(t *testing.T) {
		_, _, _, _, err := validateIssueCreateCommentAction(&CommentAction{
			Body:   "test",
			Owner:  "test",
			Repo:   "test",
			Number: "-1",
		})
		assert.Error(t, err)
	})
	t.Run("succeeds on valid comment", func(t *testing.T) {
		body, owner, repo, number, err := validateIssueCreateCommentAction(&CommentAction{
			Body:   "test-body",
			Owner:  "test-owner",
			Repo:   "test-repo",
			Number: "1",
		})
		assert.NoError(t, err)
		assert.Equal(t, "test-body", body)
		assert.Equal(t, "test-owner", owner)
		assert.Equal(t, "test-repo", repo)
		assert.Equal(t, 1, number)
	})
}

func Test_validateIssueCreateAction(t *testing.T) {
	t.Run("fails on empty owner", func(t *testing.T) {
		err := validateIssueCreateAction(&IssueCreateAction{
			Repo: "test",
		})
		assert.Error(t, err)
	})
	t.Run("fails on empty repo", func(t *testing.T) {
		err := validateIssueCreateAction(&IssueCreateAction{
			Owner: "test",
		})
		assert.Error(t, err)
	})
}

func Test_durationStringToContext(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		ctx, cancel, err := durationStringToContext("")
		require.NoError(t, err)
		t.Cleanup(cancel)
		assert.Equal(t, context.Background(), ctx)
	})

	t.Run("invalid", func(t *testing.T) {
		_, _, err := durationStringToContext("invalid")
		require.Error(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		ctx, cancel, err := durationStringToContext("1s")
		require.NoError(t, err)
		t.Cleanup(cancel)
		assert.NotEqual(t, context.Background(), ctx)
	})
}
