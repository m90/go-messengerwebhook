default: test vet

vet:
	@go vet ./...

test:
	@go test -cover ./...

.PHONY: test vet
