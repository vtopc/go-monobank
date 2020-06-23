.PHONY: test
test:
	go test `go list ./... | grep -v '/mocks'` -cover -count=1 -coverprofile=coverage.txt -covermode=count

.PHONY: deps
deps:
	go mod tidy
	go mod download
