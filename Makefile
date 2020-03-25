
.PHONY: all # Everything is accessible by user
.DEFAULT: help # Running Make will run help by defaults

help: ## Show Help
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: # Build the main package
	go build -v . 

lint: ### Runs the linter
	golangci-lint run --enable-all
  
format: ## run gofmt
	@gofmt -l -w -s .

test: ## Run Go tests
	go test -v ./...

test-coverage: ## Run Go tests with coverage
	go test -v ./... -cover	-coverprofile=coverage.txt
