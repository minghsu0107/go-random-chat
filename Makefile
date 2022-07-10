# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install

.PHONY: proto

all: build test
test:
	$(GOTEST) -gcflags=-l -v -cover -coverpkg=./... -coverprofile=cover.out ./...
build: dep
	$(GOBUILD) -ldflags="-X github.com/minghsu0107/go-random-chat/cmd.Version=v0.0.0 -w -s" -o server ./randomchat.go
dep: wire
	$(shell $(GOCMD) env GOPATH)/bin/wire ./internal/wire
wire:
	GO111MODULE=on $(GOINSTALL) github.com/google/wire/cmd/wire@v0.4.0
proto:
	protoc proto/*/*.proto --go_out=plugins=grpc:.

docker: docker-api docker-web
docker-api:
	docker build -f ./build/Dockerfile.api --build-arg VERSION=v0.0.0 -t minghsu0107/random-chat-api:kafka .
docker-web:
	docker build -f ./build/Dockerfile.web --build-arg VERSION=v0.0.0 -t minghsu0107/random-chat-web:kafka .

clean:
	$(GOCLEAN)
	rm -f server