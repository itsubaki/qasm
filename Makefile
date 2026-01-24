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
	cat testdata/grover_2x2.qasm | go run main.go -top 8

counting:
	cat testdata/counting.qasm | go run main.go -top 8

lex:
	cat testdata/bell.qasm | go run main.go -lex

parse:
	cat testdata/bell.qasm | go run main.go -parse

