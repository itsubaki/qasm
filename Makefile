SHELL := /bin/bash

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

vet:
	go vet ./...

bench:
	go test -bench . ./... --benchmem

doc:
	godoc -http=:6060
