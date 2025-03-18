build:
	@go build -o ./bin/lazyenv .

run: build
	@./bin/lazyenv

test:
	@go test ./...
