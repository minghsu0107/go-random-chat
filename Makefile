# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install

all: build test
test:
	$(GOTEST) -gcflags=-l -v -cover -coverpkg=./... -coverprofile=cover.out ./...
build: dep
	$(GOBUILD) -ldflags="-X github.com/minghsu0107/go-random-chat/cmd.Version=v0.0.0 -w -s" -o server ./randomchat.go
dep: wire
	$(shell $(GOCMD) env GOPATH)/bin/wire ./internal/wire
wire:
	GO111MODULE=on $(GOINSTALL) github.com/google/wire/cmd/wire@v0.4.0
clean:
	$(GOCLEAN)
	rm -f server