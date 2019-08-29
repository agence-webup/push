VERSION := 1.3

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-s -w"
BINARY_NAME=push
BINARY_AMD64=$(BINARY_NAME)_linux_amd64
BINARY_ARM=$(BINARY_NAME)_linux_arm

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_AMD64) -v cmd/push/main.go

publish:
	docker build -t webup/push:$(VERSION) .
	docker push webup/push:$(VERSION)
