.PHONY:
all: clean test build

.PHONY: build
build:
	mkdir -p build
	go build -o build ./...

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean
	rm -f build/*
