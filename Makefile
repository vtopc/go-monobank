.PHONY: test
test:
	go test `go list ./... | grep -v '/mocks'` -cover -count=1 -coverprofile=coverage.txt -covermode=count

.PHONY: deps
deps:
	go mod tidy
	go mod download

# linter:
GOLINT = $(GOPATH)/bin/golangci-lint
$(GOLINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.35.0

.PHONY: lint
lint: $(GOLINT)
	$(GOLINT) run
