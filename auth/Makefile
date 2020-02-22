
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/auth/auth.proto

.PHONY: build
build: proto

	CGO_ENABLED=0 GOOS=linux go build -o auth-srv *.go

.PHONY: test
test:
	go test -v ./... -cover
