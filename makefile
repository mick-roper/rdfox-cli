VERSION ?= 0.0.$(shell git rev-parse --short HEAD)-dev
GOARCH ?= amd64
GOOS ?= darwin

.PHONY: all
all: clean install test build

.PHONY: clean
clean:
	rm -rf vendor/ bin/

.PHONY: install
install:
	go mod vendor

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=$(GOOS) go build -ldflags="-s -w -X main.Version=$(VERSION)" -a -o bin/$(GOOS)/$(GOARCH)/app main.go