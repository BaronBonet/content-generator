export VERSION := $(shell ./scripts/get-version.sh)

.PHONY: clean
clean:
	@find . -type f -name 'mock_*' -delete
	@echo "Done cleaning"

.PHONY: generate-go
generate-go: clean
	go generate ./...
	@echo "Done generating go"


.PHONY: build-cli
build-cli: generate-go
	@echo "Building CLI version ${VERSION}"
	@go get ./... && \
		go build -ldflags "-X 'github.com/BaronBonet/content-generator/internal/infrastructure.Version=${VERSION}'" \
		-o content-generator cmd/cli/*go
	@chmod +x content-generator
	@echo "Done building CLI."


.PHONY: test
test: generate-go
	@go test -v ./...