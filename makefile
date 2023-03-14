VERSION ?= 0.0.$(shell git rev-parse --short HEAD)-dev

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
build: build-linux build-windows build-darwin

.PHONY: build-linux
build-linux:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.Version=$(VERSION)" -a -o bin/linux/app main.go

.PHONY: build-darwin
build-darwin:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=darwin go build -ldflags="-s -w -X main.Version=$(VERSION)" -a -o bin/darwin/app main.go

.PHONY: build-windows
build-windows:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=windows go build -ldflags="-s -w -X main.Version=$(VERSION)" -a -o bin/win/app.exe main.go