# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install

SVCS=chat match uploader user

.PHONY: proto doc

all: build test
test:
	$(GOTEST) -gcflags=-l -v -cover -coverpkg=./... -coverprofile=cover.out ./...
build: dep doc
	$(GOBUILD) -ldflags="-X github.com/minghsu0107/go-random-chat/cmd.Version=v0.0.0 -w -s" -o server ./randomchat.go

dep: wire
	$(shell $(GOCMD) env GOPATH)/bin/wire ./internal/wire
proto:
	protoc proto/*/*.proto --go_out=plugins=grpc:.
doc: swag
	for svc in $(SVCS); do \
		$(shell $(GOCMD) env GOPATH)/bin/swag init -g http.go -d pkg/$$svc -o docs/$$svc --instanceName $$svc --parseDependency --parseInternal; \
	done

wire:
	GO111MODULE=on $(GOINSTALL) github.com/google/wire/cmd/wire@v0.4.0
swag:
	GO111MODULE=on $(GOINSTALL) github.com/swaggo/swag/cmd/swag@v1.8.3

docker: docker-api docker-web
docker-api:
	@docker build -f ./build/Dockerfile.api --build-arg VERSION=v0.0.0 -t minghsu0107/random-chat-api:kafka .
docker-web:
	@docker build -f ./build/Dockerfile.web --build-arg VERSION=v0.0.0 -t minghsu0107/random-chat-web:kafka .
clean:
	$(GOCLEAN)
	rm -f server