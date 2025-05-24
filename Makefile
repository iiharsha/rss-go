.PHONY: build
build:
	@go build -o bin/rss-go .
.PHONY: run
run: build
	@./bin/rss-go
