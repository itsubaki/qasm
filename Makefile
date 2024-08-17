SHELL := /bin/bash

test:
	go test -v -cover $(shell go list ./... | grep -v /cmd | grep -v /repl | grep -v -E "qasm$$") -v -coverprofile=coverage.txt -covermode=atomic

.PHONY: repl
repl:
	go run main.go repl

lex:
	go run main.go lex -f testdata/bell.qasm
	go run main.go lex -f testdata/shor_def.qasm

parse:
	go run main.go parse -f testdata/bell.qasm
	go run main.go parse -f testdata/shor_def.qasm

.PHONY: testdata
testdata:
	go run main.go -verbose -f testdata/bell.qasm
	go run main.go -verbose -f testdata/shor.qasm

install:
	-rm ${GOPATH}/bin/qasm
	go get -u
	go mod tidy
	go install

vet:
	go vet ./...

bench:
	go test -bench . ./... --benchmem

doc:
	godoc -http=:6060
