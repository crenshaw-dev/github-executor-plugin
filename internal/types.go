package plugin

// PluginSpec represents the `plugin` block of an Argo Workflows template.
type PluginSpec struct {
	GitHub ActionSpec `json:"github,omitempty"`
}

type ActionSpec struct {
	PullRequest *PullRequestActionSpec `json:"pullRequest,omitempty"`
	Timeout     string                 `json:"timeout,omitempty"`
}

type PullRequestActionSpec struct {
	Comment *CommentAction `json:"comment,omitempty"`
}

type CommentAction struct {
	Body   string `json:"body,omitempty"`
	Owner  string `json:"owner,omitempty"`
	Repo   string `json:"repo,omitempty"`
	Number string `json:"number,omitempty"`
}
