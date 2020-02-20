SHELL = /bin/bash

ifneq ("$(wildcard .env)","")
	include .env
	export
endif

APP_NAME = dyndns

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## name: Output app name
.PHONY: name
name:
	@printf "%s" "$(APP_NAME)"

## version: Output last commit sha1
.PHONY: version
version:
	@printf "%s" "$(shell git rev-parse --short HEAD)"
