SHELL := /bin/bash

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -coverprofile=coverage.txt -covermode=atomic

repl:
	go run main.go repl

lex:
	go run main.go lex -f testdata/bell.qasm
	go run main.go lex -f testdata/teleportation.qasm
	go run main.go lex -f testdata/shor.qasm
	go run main.go lex -f testdata/gate.qasm

parse:
	go run main.go parse -f testdata/bell.qasm
	go run main.go parse -f testdata/teleportation.qasm
	go run main.go parse -f testdata/shor.qasm
	go run main.go parse -f testdata/gate.qasm

.PHONY: testdata
testdata:
	go run main.go -f testdata/bell.qasm
	go run main.go -f testdata/teleportation.qasm
	go run main.go -f testdata/shor.qasm

install:
	-rm ${GOPATH}/bin/qasm
	go mod tidy
	go install

vet:
	go vet ./...

bench:
	go test -bench . ./... --benchmem

doc:
	godoc -http=:6060
