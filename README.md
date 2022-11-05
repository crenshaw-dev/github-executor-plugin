<!-- This is an auto-generated file. DO NOT EDIT -->
# github

* Needs: >= v3.3
* Image: crenshaw-dev/github-executor-plugin:latest

This is an Argo Workflow executor plugin for interacting with GitHub.

## Example

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: github-example-
spec:
  entrypoint: main
  templates:
    - name: main
      plugin:
        github:
          # Use `issue` to create comments for PRs - the GitHub API considers PRs to be issues.
          issue:
            comment:
              body: "Hello, world!"
              number: "1"  # PR number, from the 
              owner: crenshaw-dev
              repo: github-executor-plugin
```

## Setup

## 1. Set up a GitHub personal access token

See [GitHub's instructions](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
to set up your token.

Then create a secret using that token.

```bash
# First, copy your token to the clipboard.
pbpaste > token
kubectl create secret generic github-token --from-file=token -n argo
rm token
```

## 2. Install the plugin

```bash
kubectl apply -f https://raw.githubusercontent.com/crenshaw-dev/github-executor-plugin/main/plugin.yaml
```

## 3. Run a workflow

```bash
kubectl apply -f https://raw.githubusercontent.com/crenshaw-dev/github-executor-plugin/main/github-example-workflow.yaml
```


Install:

    kubectl apply -f github-executor-plugin-configmap.yaml

Uninstall:
	
    kubectl delete cm github-executor-plugin 
