
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/upload/upload.proto

.PHONY: build
build: proto

	CGO_ENABLED=0 GOOS=linux go build -o upload-srv *.go

.PHONY: docker
docker:
	docker build . -t upload-srv:latest
