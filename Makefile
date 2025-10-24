VERSION := $(shell cat VERSION)
NAME="ufape-crawler-golang"

.PHONY: release changelog bump patch minor major

## Release Management
# Use this makefile to manage releases, changelogs, and version bumps.

release:
	git tag -a $(VERSION) -m "release: $(VERSION)"
	git-chglog -o CHANGELOG.md
	git push origin $(VERSION)
	git add CHANGELOG.md
	git commit -m "docs: atualiza changelog para $(VERSION)"
	git push

generate-tag:
	@git tag -a $(VERSION) -m "release: $(VERSION)"

changelog:
	git-chglog -o CHANGELOG.md

bump:
	@echo "Use: make [ patch | minor | major ]"

patch:
	@./scripts/bump.sh patch

minor:
	@./scripts/bump.sh minor

major:
	@./scripts/bump.sh major

## Build and Deploy
# Use this makefile to build and deploy the application for different platforms.

build:
	@make build-linux && make build-linux-static && make build-windows && make build-mac-arm64 && make build-lib-mac && make build-lib-linux && make build-lib-windows

build-snapshot:
	@$(MAKE) build NAME=$(NAME)-snapshot-$(shell date +"%Y-%m-%d")

build-snapshot-time:
	@$(MAKE) build NAME=$(NAME)-snapshot-$(shell date +"%Y-%m-%d_%H-%M-%S")

build-linux:
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o ./dist/$(NAME)-linux cmd/api/main.go

build-linux-static:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -a -installsuffix cgo -o ./dist/$(NAME)-linux-static cmd/api/main.go

build-windows:
	@GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o ./dist/$(NAME)-windows.exe cmd/api/main.go

build-mac-arm64:
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o ./dist/$(NAME)-mac-arm64 cmd/api/main.go

build-lib-mac:
	@go build -buildmode=c-shared -ldflags="-X 'main.Version=$(VERSION)'" -o ./dist/$(NAME)-mac-arm64.dylib cmd/api/main.go

build-lib-linux:
	@go build -buildmode=c-shared -ldflags="-X 'main.Version=$(VERSION)'" -o ./dist/$(NAME)-linux.so cmd/api/main.go

build-lib-windows:
	@go build -buildmode=c-shared -ldflags="-X 'main.Version=$(VERSION)'" -o ./dist/$(NAME)-windows.dll cmd/api/main.go

## Development and Testing
# Use this makefile to run the application, manage Docker containers, and handle application.

run:
	@go run cmd/api/main.go

docker-up:
	@docker-compose -f ./docker-compose.yaml up -d --build

docker-down:
	@docker-compose -f ./docker-compose.yaml down

docker-clean:
	@docker system prune --all --volumes --force

docker-create-image:
	@echo "Building Docker image nettojulio/$(NAME):$(VERSION)"
	@docker build -t nettojulio/$(NAME):$(VERSION) .
	@docker tag nettojulio/$(NAME):$(VERSION) nettojulio/$(NAME):latest

docker-push-image:
	@make docker-create-image
	@echo "Pushing Docker image nettojulio/$(NAME):$(VERSION) to Docker Hub"
	@docker push nettojulio/ufape-crawler-golang:$(VERSION)
	@docker push nettojulio/ufape-crawler-golang:latest

## Dependency Management
# Use this makefile to manage Go dependencies, check for updates, and clean up unused dependencies.

check-updates:
	@go list -u -m all

update-deps:
	@go get -u

cleanup-deps:
	@go mod tidy

update-go-version:
	@GO_VERSION=$$(go env GOVERSION | sed -E 's/[^0-9.]//g') && \
	echo "Atualizando para Go version $$GO_VERSION" && \
	go mod edit -go=$$GO_VERSION

cleanup-cache:
	@go clean -cache -modcache -i -r

## Testing and Linting
test:
	@go test -v ./...

lint:
	@go vet ./...
	@golangci-lint run --timeout 5m

fmt:
	@go fmt ./...
	@gofumpt -l -w .
	@goimports -w .

vet:
	@go vet ./...
	@golangci-lint run --timeout 5m

vet-fix:
	@go vet -fix ./...
	@golangci-lint run --fix --timeout 5m

cover-profile:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -html=coverage.out

cover-aggregate:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

## Swagger Documentation

generate-swagger:
	@swag init -g cmd/api/main.go
	@echo "Swagger documentation generated in docs/swagger"
