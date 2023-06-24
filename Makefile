export VERSION := $(shell ./scripts/get-version.sh)

OUT_DIR=out
ARCH=arm64
OS=linux

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

.PHONY: build-go-aws
build-go-aws: generate-go
	@echo "Building content generator for AWS version ${VERSION}"
	@go get ./... && \
		GOOS=${OS} GOARCH=${ARCH} go build -ldflags \
		 "-X 'github.com/BaronBonet/content-generator/internal/infrastructure.Version=${VERSION}'" \
		 -o ${OUT_DIR}/handler/bootstrap cmd/aws_lambda/*go
	@zip -jrm "./${OUT_DIR}/handler/main.zip" "./${OUT_DIR}/handler/"*
	@echo "Done building for AWS."

.PHONY: test
test: generate-go
	@go test -v ./...


.PHONY: tf-plan
tf-plan:
	@cd ./terraform && terraform plan -var-file=variables.tfvars

.PHONY: tf-apply
tf-apply:
	@cd ./terraform && terraform apply -var-file=variables.tfvars
