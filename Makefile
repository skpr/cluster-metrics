#!/usr/bin/make -f

export CGO_ENABLED=0
export GO111MODULE=on

REGISTRY=skpr/cluster-metrics
OUTPUT=bin/cluster-metrics
VERSION=$(shell git describe --tags --always)

default: lint test build

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run go fmt against code
fmt:
	go fmt ./...

vet:
	go vet ./...

# Run tests with coverage reporting.
test:
	gotestsum -- -coverprofile=cover.out ./...

build:
	GOOS=linux go build -o ${OUTPUT} main.go

docker:
	docker build -t ${REGISTRY}:${VERSION} -t ${REGISTRY}:latest .

push:
	docker push ${REGISTRY}:${VERSION}
	docker push ${REGISTRY}:latest

.PHONY: *
