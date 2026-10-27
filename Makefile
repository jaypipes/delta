# Top-level Makefile containing make targets that automate build, release,
# deployment and local testing tasks.

ROOT_DIR ?= $(shell pwd)

GIT_VERSION ?= $(shell git describe --tags --always --dirty || echo "unknown")
GIT_COMMIT ?= $(shell git rev-parse HEAD)

include .mk/go.mk
include .mk/test.mk
include .mk/help.mk
