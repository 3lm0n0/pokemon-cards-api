NAME=cards
BINARY_NAME=cards
MAIN_PACKAGE_PATH=./cmd/api/v1/main.go
BINARY_PATH=./bin/${BINARY_NAME}
LDFLAGS="all=-w -s"
OS=linux
ARCH_X86_64=x86_64
ARCH_ARM64=arm64
ARCH_AMD64=amd64

.PHONY:= hello init setup build run test clean
.DEFAULT_GOAL:= setup build run

# ==================================================================================== #
# HELP
# ==================================================================================== #
## help: prints this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# CHECKHEALTH
# ==================================================================================== #
## ping: responds with pong to ensure makefile health
.PHONY: ping
ping:
	echo "pong"

# ==================================================================================== #
# SETUP
# ==================================================================================== #
## init: initialize the application
.PHONY: init
init:
	@echo "=> Go module ${BINARY_NAME} initializing"
	@go mod init '${BINARY_NAME}'
	@echo "=> Go module initialized"

## setup: sets-up the application
.PHONY: setup
setup:
	@echo "=> Stetting microservice"
	@export GOSUMDB=off
	@go mod tidy
	@go mod download
	@echo "=> Setup completed"

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
## build: builds the application
.PHONY: build
build:
	@echo "=> Building microservice"
	@env GOOS=${OS} GOARCH=${ARCH_AMD64} go build -ldflags=${LDFLAGS} -o ${BINARY_PATH} ${MAIN_PACKAGE_PATH}

	@echo "=> Building completed"

## run: runs the application
.PHONY: run
run:
	./bin/${BINARY_NAME}

## test: tests the application
.PHONY: test
test:
	@go test -v ./...

## test-race: tests races in the application
.PHONY: test-race
test-race:
	@go test -race -buildvcs -vet=off ./...

## test-cover: test application coverage
.PHONY: test-cover
test-cover:
	@go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./... 
	@go tool cover -html=/tmp/coverage.out

## clean: cleans all binaries, objects and sum ...
.PHONY: clean
clean:
	@echo "Cleaning up all binaries, objects and sum ..."
	@go clean
	@rm -rvf *.o ./bin/${BINARY_NAME} go.sum
	@echo "Cleaning completed"

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #
##production-deploy: deploys into production enviroment
.PHONY: production-deploy
production-deploy: build
	@echo "=> Starting deploy"
	@serverless deploy --stage prod
	@echo "=> Deploy finished"