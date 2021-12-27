NAME ?= "content-plus-totp"

.PHONY: test
test:
	go test -race -v $$(go list ./...)

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: test-all
test-all: test lint

.PHONY: compile
compile:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o dist/${NAME} main.go
