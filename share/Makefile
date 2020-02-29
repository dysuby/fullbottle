
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/share/share.proto

.PHONY: build
build: proto

	CGO_ENABLED=0 GOOS=linux go build -o share-srv *.go

.PHONY: docker
docker:
	docker build . -t share-srv:latest
