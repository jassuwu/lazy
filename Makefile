build:
	@go build -o ./bin/lazy .

run: build
	@./bin/lazy

test:
	@go test ./...
