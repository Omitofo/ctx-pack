.PHONY: build run version clean

build:
	go build -o ctx-pack ./cmd/ctx-pack

run: build
	./ctx-pack

version: build
	./ctx-pack --version

clean:
	rm -f ctx-pack
