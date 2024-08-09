PACKNAME={{ .ServiceName }}
BUILD_FLAGS=-X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`'
BUILD_WIN=@env GOOS=windows GOARCH=amd64 go build -o $(PACKNAME).exe
BUILD_LINUX=@env GOOS=linux GOARCH=amd64 go build -o $(PACKNAME)
BUILD_MAC=@env GOOS=darwin GOARCH=amd64 go build -o $(PACKNAME)-darwin

release: go-package go-release go-bin
release-candidate: go-package go-release-candidate
binary: go-package go-build

go-release: release-nix release-win release-mac

go-build:
	@echo "  >  Building binary ..."
	@go build -o $(PACKNAME) -ldflags="$(BUILD_FLAGS)"

go-package: go-tidy go-test
	@echo "  >  Packaging static files..."

go-test:
	@echo "  >  Validating code ..."
	@go vet ./...
	@go test ./...

go-tidy:
	@echo "  >  Tidying go.mod ..."
	@go mod tidy

go-test-cov:
	@echo "Running tests and generating coverage output ..."
	@go test ./... -coverprofile coverage.out -covermode count
	@sleep 2 # Sleeping to allow for coverage.out file to get generated
	@echo "Current test coverage : $(shell go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+') %"

go-release-candidate: go-tidy go-test
	@echo "  >  Building release candidate for Linux..."
	$(BUILD_LINUX) -ldflags="$(BUILD_FLAGS) -X 'main.VersionPostfix=nix-rc'"
	@echo "  >  Building release candidate for Windows..."
	$(BUILD_WIN) -ldflags="$(BUILD_FLAGS) -X 'main.VersionPostfix=win-rc'"
	@echo "  >  Building release for Darwin..."
	$(BUILD_MAC) -ldflags="$(BUILD_FLAGS) -X 'main.VersionPostfix=darwin-rc'"

release-nix:
	@echo "  >  Building release for Linux..."
	$(BUILD_LINUX) -ldflags="$(BUILD_FLAGS) -X 'main.VersionPostfix=linux'"

release-win:
	@echo "  >  Building release for Windows..."
	$(BUILD_WIN) -ldflags="$(BUILD_FLAGS) -X 'main.VersionPostfix=windows'"

release-mac:
	@echo "  >  Building release for Darwin..."
	$(BUILD_MAC) -ldflags="$(BUILD_FLAGS) -X 'main.VersionPostfix=darwin'"

go-bin:
	@mv $(PACKNAME)* ~/privateer/bin