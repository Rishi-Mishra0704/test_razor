
.PHONY: build run
clean:
	@rm -rf bin

build: clean
	@templ generate
	@go build -o bin/test_razor ./cmd/...

run: build
	@./bin/test_razor