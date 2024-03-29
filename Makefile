# Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
# Use of this source code is governed by
# MIT license that can be found in the LICENSE file.


.PHONY: all clean nuke
grep=--include=*.go  --exclude=*.txt
TOP = $(shell pwd)
SRC := $(shell find $(TOP) -iname "*.go" -exec grep -L  "Copyright (c) 2021" {} \;)

all: build test
	go vet ./...
	staticcheck ./...
	make todo

build: require
	rm -rf out/*.go
	touch Builder/GoCodeTemplate.go
	@echo "package builder\n var goCodeTemplateStr string=\`" > Builder/GoCodeTemplate.go
	@cat  Builder/GoCodeTemplate.go Builder/goCode.templ > Builder/GoCodeTemplate.go.tmp
	echo "\`" >> Builder/GoCodeTemplate.go.tmp
	mv Builder/GoCodeTemplate.go.tmp Builder/GoCodeTemplate.go
	touch Builder/GoObjectTemplate.go
	@echo "package builder\n var goObjectTemplateStr string=\`" > Builder/GoObjectTemplate.go
	@cat  Builder/GoObjectTemplate.go Builder/goObject.templ > Builder/GoObjectTemplate.go.tmp
	echo "\`" >> Builder/GoObjectTemplate.go.tmp
	mv Builder/GoObjectTemplate.go.tmp Builder/GoObjectTemplate.go
	go fmt ./...
	go build -o bin/yaccgo ./yaccgo/*.go

require:
	@go version >/dev/null 2>&1 || { echo >&2 "go is required but not installed.  Aborting."; exit 1; }
	@staticcheck --version >/dev/null 2>&1 || { go install honnef.co/go/tools/cmd/staticcheck@2020.2.1; }
todo:
	@grep -nr $(grep) ^[[:space:]]*_[[:space:]]*=[[:space:]][[:alpha:]][[:alnum:]]* * || true
	@grep -nr $(grep) TODO * || true
	@grep -nr $(grep) BUG * || true
	@grep -nr $(grep) [^[:alpha:]]println * || true
test: require
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
header:
	@for f in $(SRC); do \
        cat $(TOP)/head $$f > $$f.tmp && mv $$f.tmp $$f; \
    done