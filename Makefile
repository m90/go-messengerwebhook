default: test

test:
	@go test -race -v -cover ./...

.PHONY: test
