.PHONY: test
test:
	go test -timeout 30s ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "#####"
	@echo "If running on WSL, open the report with: "
	@echo "explorer.exe coverage.html"

.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: generate
generate:
	go generate ./...
