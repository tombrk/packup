# VERSION := $(shell git describe --tags --dirty --always)

.PHONY: ui
ui:
	cd ui && yarn build
