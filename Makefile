default: test vet

vet:
	@go vet ./...

test:
	@go test -race -v -cover ./...

.PHONY: test vet
