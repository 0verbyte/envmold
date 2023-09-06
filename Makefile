.PHONY: build
build:
	go build -o bin/mold .

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run -v
