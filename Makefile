.PHONY: build
build:
	docker build --load -t crenshawdotdev/github-executor-plugin:latest .

.PHONY: manifests
manifests:
	argo executor-plugin build .
