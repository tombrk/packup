.DEFAULT_GOAL: dev

VERSION := $(shell git describe --tags --dirty --always)

# run dev server
.PHONY: dev
dev: pkged.go
	CGO_ENABLED=0 go build -ldflags=${LDFLAGS} .

LDFLAGS := '-s -w -extldflags "-static" -X main.Version=${VERSION}'
static: pkged.go
	CGO_ENABLED=0 go build -o packup -ldflags=${LDFLAGS} .

# pack ui into Go
pkged.go: ui/build/.gitkeep
	pkger

# rebuild ui
.PHONY: react
react:
	cd ui && yarn build
