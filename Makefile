
PACKAGE = github.com/mikattack/golang-rest-stub
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`

.PHONY: vendor check fmt test coverage
.DEFAULT_GOAL := help

server:
	go build "${PACKAGE}/cmd/rest-stub"

logger:
	go build "${PACKAGE}/cmd/rest-stub-log"

check: test fmt

vendor:
	go get github.com/chilts/sid
	go get github.com/go-kit/kit/log
	go get github.com/oklog/run
	go get github.com/stretchr/testify

test: vendor
	go test

fmt:
	gofmt -l *.go

coverage:
	# echo "mode: count" > coverage-all.out
	go test -coverprofile=coverage.out -covermode=count $(pkg);
	# tail -n +2 coverage.out >> coverage-all.out;
	go tool cover -html=coverage.out
