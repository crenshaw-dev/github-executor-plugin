.PHONY: build
build:
	docker build --load -t crenshaw-dev/github-executor-plugin:latest .

.PHONY: manifests
manifests:
	argo executor-plugin build ./manifests
