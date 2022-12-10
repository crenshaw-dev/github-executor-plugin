package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_validateIssueAction(t *testing.T) {
	t.Run("fails on nil action", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(nil)
		assert.Error(t, err)
	})
	t.Run("fails on nil comment", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(&IssueActionSpec{})
		assert.Error(t, err)
	})
	t.Run("fails on empty comment body", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(&IssueActionSpec{
			Comment: &CommentAction{},
		})
		assert.Error(t, err)
	})
	t.Run("fails on empty owner", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(&IssueActionSpec{
			Comment: &CommentAction{
				Body: "test",
			},
		})
		assert.Error(t, err)
	})
	t.Run("fails on empty repo", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(&IssueActionSpec{
			Comment: &CommentAction{
				Body:  "test",
				Owner: "test",
			},
		})
		assert.Error(t, err)
	})
	t.Run("fails on empty number", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(&IssueActionSpec{
			Comment: &CommentAction{
				Body:  "test",
				Owner: "test",
				Repo:  "test",
			},
		})
		assert.Error(t, err)
	})
	t.Run("fails on invalid number", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(&IssueActionSpec{
			Comment: &CommentAction{
				Body:   "test",
				Owner:  "test",
				Repo:   "test",
				Number: "test",
			},
		})
		assert.Error(t, err)
	})
	t.Run("fails on negative number", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(&IssueActionSpec{
			Comment: &CommentAction{
				Body:   "test",
				Owner:  "test",
				Repo:   "test",
				Number: "-1",
			},
		})
		assert.Error(t, err)
	})
	t.Run("succeeds on valid comment", func(t *testing.T) {
		body, owner, repo, number, err := validateIssueAction(&IssueActionSpec{
			Comment: &CommentAction{
				Body:   "test-body",
				Owner:  "test-owner",
				Repo:   "test-repo",
				Number: "1",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, "test-body", body)
		assert.Equal(t, "test-owner", owner)
		assert.Equal(t, "test-repo", repo)
		assert.Equal(t, 1, number)
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
