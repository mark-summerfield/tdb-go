#!/bin/bash
clc -s -e tdb_test.go db1_test.go eg bin
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st
