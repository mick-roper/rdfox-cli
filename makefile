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
build: build-linux build-windows

.PHONY: build-linux
build-linux:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -o bin/app main.go

.PHONY: build-windows
build-windows:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=windows go build -ldflags="-s -w" -a -o bin/app.exe main.go