
GOPATH:=$(shell go env GOPATH)


.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o api-web *.go

.PHONY: test
test:
	go test -v ./... -cover
