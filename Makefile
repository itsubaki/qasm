SHELL := /bin/bash

antlr:
	# https://github.com/openqasm/openqasm/blob/main/source/grammar
	pip install antlr4-tools
	antlr4 -Dlanguage=Go -visitor -o ./gen/parser -package parser qasm3Lexer.g4 qasm3Parser.g4

lex:
	go run cmd/lex/main.go < _testdata/bell.qasm

parse:
	go run cmd/parse/main.go < _testdata/bell.qasm

test:
	go test -v -cover $(shell go list ./... | grep -v /cmd ) -v -coverprofile=coverage.txt -covermode=atomic
