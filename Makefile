build:
	@go build -o ./bin/gohttp

run: build
	@./bin/gohttp

image:
	docker build -t gohttp .