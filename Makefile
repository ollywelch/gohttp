build:
	@go build -o ./bin/gohttp

run: build
	@./bin/gohttp