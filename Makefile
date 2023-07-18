BUILD_CHANNEL?=local
OS=$(shell uname)
VERSION=v1.12.0
GIT_REVISION = $(shell git rev-parse HEAD | tr -d '\n')
TAG_VERSION?=$(shell git tag --points-at | sort -Vr | head -n1)
CGO_LDFLAGS=""
GO_BUILD_LDFLAGS = -ldflags "-extldflags=-static -X 'main.Version=${TAG_VERSION}' -X 'main.GitRevision=${GIT_REVISION}'"

.PHONY: default
default: build-module

.PHONY: goformat
goformat:
	go install golang.org/x/tools/cmd/goimports
	gofmt -s -w .
	goimports -w -local=go.viam.com/utils `go list -f '{{.Dir}}' ./... | grep -Ev "proto"`

.PHONY: lint
lint: goformat
	go install github.com/edaniels/golinters/cmd/combined
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	go install github.com/polyfloyd/go-errorlint
	go list -f '{{.Dir}}' ./... | grep -v gen | xargs go vet -vettool=`go env GOPATH`/bin/combined
	go list -f '{{.Dir}}' ./... | grep -v gen | xargs `go env GOPATH`/bin/go-errorlint -errorf
	go list -f '{{.Dir}}' ./... | grep -v gen | xargs go run github.com/golangci/golangci-lint/cmd/golangci-lint run -v --config=./etc/.golangci.yaml

.PHONY: test
test:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...


.PHONY: build
build: 
	mkdir -p bin && CGO_ENABLED=0 CGO_LDFLAGS=${CGO_LDFLAGS} go build $(GO_BUILD_LDFLAGS) -o bin/module module/main.go

.PHONY: package
package: build
	mkdir -p package && tar -czf ./package/module.tar.gz ./bin/module


.PHONY: clean
clean: 
	rm -rf bin

