SHELL := /bin/bash

test:
	go test -v -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -coverprofile=coverage.txt -covermode=atomic

.PHONY: repl
repl:
	go run main.go repl

lex:
	go run main.go lex -f _testdata/bell.qasm
	go run main.go lex -f _testdata/shor_def.qasm

parse:
	go run main.go parse -f _testdata/bell.qasm
	go run main.go parse -f _testdata/shor_def.qasm

testdata:
	go run main.go -verbose -f _testdata/bell.qasm
	go run main.go -verbose -f _testdata/shor.qasm

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
