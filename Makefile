build:
	@go build -o ./bin/lazyenv ./cmd

run: build
	@./bin/lazyenv

test:
	@go test ./...
