APP_NAME=go-scoring-api
COVERAGE_FILE=coverage.out

test:
	@echo "🧪 Ejecutando tests..."
	go test ./... -v

test-cover:
	@echo "🧪 Ejecutando tests con cobertura..."
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

test-cover-html:
	@echo "🧪 Ejecutando tests con cobertura (HTML)..."
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE)

clean:
	@rm -f $(COVERAGE_FILE)

.PHONY: test test-cover test-cover-html clean
