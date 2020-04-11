.PHONY: test
test:
	go test `go list ./... | grep -v '/mocks'` -cover -count=1

.PHONY: deps
deps:
	go mod tidy
	go mod vendor
