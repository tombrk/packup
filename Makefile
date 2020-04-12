.DEFAULT_GOAL: dev

# run dev server
.PHONY: dev
dev: pkged.go
	go run .

# pack ui into Go
pkged.go: ui/build/.gitkeep
	pkger

# rebuild ui
.PHONY: react
react:
	cd ui && yarn build
