setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run ./cli

lint-fix:
	golangci-lint run ./cli --fix