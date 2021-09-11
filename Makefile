
# ==================================================================================== # 
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== # 
# DEVELOPMENT
# ==================================================================================== #

## run: run main application
.PHONY: run
run:
	go run ./src --displayflags 

## run/bin: run the cmd/api application
.PHONY: run/bin
run/bin:
	./bin/mock-json

# ==================================================================================== # 
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor, dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	
	@echo '==> Formatting code...'
	go fmt ./...
	@echo '==> Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo '==> Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo '==> Tidying and verifying module dependencies...' 
	go mod tidy
	go mod verify
	@echo '==> Vendor dependencies...'
	go mod vendor


## cover: roda os testes com cover 
.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ==================================================================================== # 
# BUILD
# ==================================================================================== #

current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime="${current_time}"" -X main.version=${git_description}'

## build: build the main application.
.PHONY: build
build:
	@echo '==> Building mock-json...'
	@echo ${current_time} 
	@echo ${git_description}
	go build -ldflags=${linker_flags} -o=./bin/mock-json ./src
	@echo '==> Building main to linux...'
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/mock-json ./src

