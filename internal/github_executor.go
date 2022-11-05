package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
	"github.com/google/go-github/github"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

type GitHubExecutor struct {
	client *github.Client
}

func NewGitHubExecutor(client *github.Client) GitHubExecutor {
	return GitHubExecutor{client: client}
}

func (e *GitHubExecutor) Execute(args executor.ExecuteTemplateArgs) executor.ExecuteTemplateReply {
	pluginJSON, err := args.Template.Plugin.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("failed to marshal plugin to JSON from workflow spec: %w", err)
		log.Println(err.Error())
		return errorResponse(err)
	}

	plugin := &PluginSpec{}
	err = json.Unmarshal(pluginJSON, plugin)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal plugin JSON to plugin struct: %w", err)
		log.Println(err.Error())
		return errorResponse(err)
	}

	output, err := e.runAction(plugin)
	if err != nil {
		return failedResponse(wfv1.Progress(fmt.Sprintf("0/1")), fmt.Errorf("action failed: %w", err))
	}

	outPtr := &output

	return executor.ExecuteTemplateReply{
		Node: &wfv1.NodeResult{
			Phase:    wfv1.NodeSucceeded,
			Message:  "Action completed",
			Progress: "1/1",
			Outputs: &wfv1.Outputs{
				Result: outPtr,
			},
		},
	}
}

func (e *GitHubExecutor) runAction(plugin *PluginSpec) (string, error) {
	ctx, cancel, err := durationStringToContext(plugin.GitHub.Timeout)
	if err != nil {
		return "", fmt.Errorf("failed to parse timeout: %w", err)
	}
	defer cancel()
	body, owner, repo, number, err := validatePullRequestAction(plugin.GitHub.Issue)
	if err != nil {
		return "", fmt.Errorf("invalid pull request action: %w", err)
	}
	_, response, err := e.client.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{
		Body: &body,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create pull request comment: %w", err)
	}
	if response.StatusCode != 201 {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}
		return "", fmt.Errorf("failed to create pull request comment: %s", string(responseBody))
	}
	return "", nil
}

func validatePullRequestAction(action *IssueActionSpec) (body, owner, repo string, number int, err error) {
	if action == nil {
		return "", "", "", -1, fmt.Errorf("the only available action for the GitHub plugin is 'pullRequest'")
	}
	if action.Comment == nil {
		return "", "", "", -1, fmt.Errorf("the only available action for pull requests is `comment`")
	}
	if action.Comment.Body == "" {
		return "", "", "", -1, fmt.Errorf("the pull request comment body is required")
	}
	if action.Comment.Owner == "" {
		return "", "", "", -1, fmt.Errorf("the pull request owner is required")
	}
	if action.Comment.Repo == "" {
		return "", "", "", -1, fmt.Errorf("the pull request repo is required")
	}
	if action.Comment.Number == "" {
		return "", "", "", -1, fmt.Errorf("the pull request number is required")
	}
	number, err = strconv.Atoi(action.Comment.Number)
	if err != nil {
		return "", "", "", -1, fmt.Errorf("the pull request number must be an integer")
	}
	if number < 0 {
		return "", "", "", -1, fmt.Errorf("the pull request number must be greater than or equal to 0")
	}
	return action.Comment.Body, action.Comment.Owner, action.Comment.Repo, number, nil
}

// durationStringToContext parses a duration string and returns a context and cancel function. If timeout is empty, the
// context is context.Background().
func durationStringToContext(timeout string) (ctx context.Context, cancel func(), err error) {
	ctx = context.Background()
	cancel = func() {}
	if timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse timeout: %w", err)
		}
		ctx, cancel = context.WithTimeout(ctx, duration)
	}
	return ctx, cancel, nil
}

func errorResponse(err error) executor.ExecuteTemplateReply {
	return executor.ExecuteTemplateReply{
		Node: &wfv1.NodeResult{
			Phase:    wfv1.NodeError,
			Message:  err.Error(),
			Progress: wfv1.ProgressZero,
		},
	}
}

func failedResponse(progress wfv1.Progress, err error) executor.ExecuteTemplateReply {
	return executor.ExecuteTemplateReply{
		Node: &wfv1.NodeResult{
			Phase:    wfv1.NodeFailed,
			Message:  err.Error(),
			Progress: progress,
		},
	}
}
