VERSION := $(shell git describe --tags --dirty --always)
LDFLAGS := '-s -w -extldflags "-static" -X main.Version=${VERSION}'
GO := CGO_ENABLED=0 go build -trimpath -ldflags=${LDFLAGS}

.PHONY: server
server: ## packup server (UI and metrics)
	$(GO) -o packup-server .

.PHONY: agent
agent: ## packup agent (run restic on schedule)
	$(GO) -o packup-agent ./agent

.PHONY: ui
ui: ## build web-ui using bun
	cd ui && bun run build

.PHONY: docker
docker: ## build images
	docker build -t shorez/packup .
	docker build -t shorez/packup-agent --target=agent .

PLATFORMS:=linux/arm64,linux/amd64
docker-cross: ## build and push images for amd64,arm64
	docker buildx build -t shorez/packup --platform=$(PLATFORMS) --push .
	docker buildx build -t shorez/packup-agent --platform=$(PLATFORMS) --target=agent --push .

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST) | sort
.DEFAULT_GOAL=help
.PHONY=help
