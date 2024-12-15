SHELL := /bin/bash

antlr:
	# https://github.com/openqasm/openqasm/blob/main/source/grammar
	pip install antlr4-tools
	antlr4 -Dlanguage=Go -visitor -o ./parser -package parser qasm3Lexer.g4 
	antlr4 -Dlanguage=Go -visitor -o ./parser -package parser qasm3Parser.g4 
