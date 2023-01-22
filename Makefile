.PHONY:
all: clean test build

.PHONY: build
build:
	mkdir -p build
	go build -o build ./...

.PHONY: test
test:
	go test ./...

.PHONY: snapshot
snapshot:
	goreleaser release --snapshot --rm-dist

.PHONY: release
release:
	goreleaser release

.PHONY: clean
clean:
	go clean
	rm -f build/*
