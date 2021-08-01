SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

install:
	-rm ${GOPATH}/bin/qasm
	go mod tidy
	go install -ldflags "${LDFLAGS}"

vet:
	go vet ./...

bench:
	go test -bench . ./... --benchmem

doc:
	godoc -http=:6060
