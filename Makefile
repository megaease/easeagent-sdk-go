.PHONY: build

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))

build:
	CGO_ENABLED=0 go build -v -o ${MKFILE_DIR}bin/meshdemo ${MKFILE_DIR}example/mesh/main.go

build_docker:
	docker buildx build --platform linux/amd64 --load -t megaease/meshdemo:latest -f ./example/mesh/Dockerfile .
	docker tag megaease/meshdemo:latest megaease/meshdemo:canary_version
