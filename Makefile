export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1
GOOS=$(shell go env GOOS)
VERSION=$(shell git describe --tags --always)
REVISION=$(shell git rev-parse HEAD)$(shell if ! git diff --no-ext-diff --quiet --exit-code; then echo .m; fi)

dep:
	go mod tidy
	go get -d -u golang.org/x/tools/cmd/goimports

build:
	CGO_ENABLED=0 GOOS=${GOOS} go build -ldflags "-X=main.Revision=${REVISION} -X=main.Version=${VERSION}" -o httplogger ./cmd/main.go

lint:
	revive --config=revive.toml --formatter=unix ./...

fmt:
	go fmt ./...
	goimports -local andboson   -w .

all: build

