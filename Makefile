.PHONY: build
build:
	go build -v ./cmd/apiserver

.DEFFAULT_GOAL := build
