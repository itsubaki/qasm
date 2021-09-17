SHELL := /bin/bash

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -coverprofile=coverage.txt -covermode=atomic

repl:
	go run main.go repl

lex:
	go run main.go lex -f testdata/bell.qasm
	go run main.go lex -f testdata/bell_gate.qasm
	go run main.go lex -f testdata/shor.qasm
	go run main.go lex -f testdata/shor_def.qasm

parse:
	go run main.go parse -f testdata/bell.qasm
	go run main.go parse -f testdata/bell_gate.qasm
	go run main.go parse -f testdata/shor.qasm
	go run main.go parse -f testdata/shor_def.qasm

.PHONY: testdata
testdata:
	go run main.go -v -f testdata/bell.qasm
	go run main.go -v -f testdata/bell_ctrl.qasm
	go run main.go -v -f testdata/bell_gate.qasm
	go run main.go -v -f testdata/shor.qasm
	go run main.go -v -f testdata/mod.qasm
	go run main.go -v -f testdata/mod_negc.qasm

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
