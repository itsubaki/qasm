SHELL := /bin/bash

test:
	go test -v -cover $(shell go list ./... | grep -v /cmd | grep -v /gen | grep -v -E "qasm$$") -v -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run

update:
	GOPROXY=direct go get github.com/itsubaki/q@HEAD
	go get -u ./...
	go mod tidy

antlr:
	curl -s -O https://raw.githubusercontent.com/openqasm/openqasm/refs/heads/main/source/grammar/qasm3Lexer.g4
	curl -s -O https://raw.githubusercontent.com/openqasm/openqasm/refs/heads/main/source/grammar/qasm3Parser.g4
	pip install antlr4-tools
	antlr4 -Dlanguage=Go -visitor -o ./gen/parser -package parser qasm3Lexer.g4 qasm3Parser.g4

shor:
	cat testdata/shor15.qasm | go run main.go -top 8

grover:
	cat testdata/grover.qasm | go run main.go -top 8

counting:
	cat testdata/quantum_counting.qasm | go run main.go -top 16

lex:
	cat testdata/bell.qasm | go run main.go -lex

parse:
	cat testdata/bell.qasm | go run main.go -parse

validate:
	cat testdata/qft.qasm | go run main.go -validate

.PHONY: svg
svg:
	go run main.go -svg < testdata/svg/barrier.qasm               > testdata/svg/barrier.svg
	go run main.go -svg < testdata/svg/bell.qasm                  > testdata/svg/bell.svg
	go run main.go -svg < testdata/svg/ghz.qasm                   > testdata/svg/ghz.svg
	go run main.go -svg < testdata/svg/inv_qft.qasm               > testdata/svg/inv_qft.svg
	go run main.go -svg < testdata/svg/quantum_teleportation.qasm > testdata/svg/quantum_teleportation.svg
	go run main.go -svg < testdata/svg/shor15.qasm                > testdata/svg/shor15.svg
	go run main.go -svg < testdata/bell.qasm                      > testdata/bell.svg
	go run main.go -svg < testdata/deutsch_jozsa_balanced.qasm    > testdata/deutsch_jozsa_balanced.svg
	go run main.go -svg < testdata/deutsch_jozsa_constant.qasm    > testdata/deutsch_jozsa_constant.svg
	go run main.go -svg < testdata/error_correction.qasm          > testdata/error_correction.svg
	go run main.go -svg < testdata/grover.qasm                    > testdata/grover.svg
	go run main.go -svg < testdata/qft.qasm                       > testdata/qft.svg
	go run main.go -svg < testdata/qsp.qasm                       > testdata/qsp.svg
	go run main.go -svg < testdata/quantum_teleportation.qasm     > testdata/quantum_teleportation.svg
	go run main.go -svg < testdata/shor15.qasm                    > testdata/shor15.svg
