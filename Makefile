.PHONY: build test vet fmt lint clean

build:
	mkdir -p bin
	go build -o bin/sqlschemaview ./cmd/sqlschemaview

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w ./cmd ./internal ./integration

lint: vet
	test -z "$$(gofmt -l ./cmd ./internal ./integration)"

clean:
	rm -rf bin coverage.out
