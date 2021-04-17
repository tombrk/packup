.DEFAULT_GOAL: dev

VERSION := $(shell git describe --tags --dirty --always)

# run dev server
.PHONY: dev
dev:
	CGO_ENABLED=0 go build -ldflags=${LDFLAGS} .

LDFLAGS := '-s -w -extldflags "-static" -X main.Version=${VERSION}'
static:
	CGO_ENABLED=0 go build -o packup -ldflags=${LDFLAGS} .

# rebuild ui
.PHONY: react
react:
	cd ui && yarn build
