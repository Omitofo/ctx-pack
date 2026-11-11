.PHONY: build run version test clean

build:
	go build -o ctx-pack ./cmd/ctx-pack

run: build
	./ctx-pack

version: build
	./ctx-pack --version

test:
	go test ./...

clean:
	rm -f ctx-pack
