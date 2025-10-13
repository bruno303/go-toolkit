.PHONY: test
test:
	go test -timeout 30s ./...

.PHONY: deps
deps:
	go mod tidy
	go mod vendor
