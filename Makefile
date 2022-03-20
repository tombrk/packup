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
ui: ## React based web interface (bundled into server if built beforehand)
	cd ui && yarn build

.PHONY: docker ## Build docker images
docker:
	docker build -t shorez/packup .
	docker build -t shorez/packup-agent --target=agent .

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST) | sort
.DEFAULT_GOAL=help
.PHONY=help
