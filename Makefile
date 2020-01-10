GO ?= go

vet:
	${GO} vet ./... 

build:
	${GO} build

test: 
	${GO} test ./...	

.PHONY: build
