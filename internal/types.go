package plugin

import "github.com/google/go-github/github"

// PluginSpec represents the `plugin` block of an Argo Workflows template.
type PluginSpec struct {
	GitHub *ActionSpec `json:"github,omitempty"`
}

type ActionSpec struct {
	Issue   *IssueActionSpec `json:"issue,omitempty"`
	Timeout string           `json:"timeout,omitempty"`
}

type IssueActionSpec struct {
	Comment *IssueCommentAction `json:"comment,omitempty"`
	Create  *IssueCreateAction  `json:"create,omitempty"`
}

type IssueCommentAction struct {
	Body   string `json:"body,omitempty"`
	Owner  string `json:"owner,omitempty"`
	Repo   string `json:"repo,omitempty"`
	Number string `json:"number,omitempty"`
}

type IssueCreateAction struct {
	Owner   string               `json:"owner,omitempty"`
	Repo    string               `json:"repo,omitempty"`
	Request *github.IssueRequest `json:"-"`
}
