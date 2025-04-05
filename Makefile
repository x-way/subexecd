.PHONY: all build format lint clean

all: format lint build

build:
	go build .
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o subexecd.aarch64 .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o subexecd.amd64 .

format:
	go install mvdan.cc/gofumpt@latest
	gofumpt -w .

lint:
	go vet -v ./...
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -checks all ./...

clean:
	rm -f subexecd subexecd.aarch64 subexecd.amd64
