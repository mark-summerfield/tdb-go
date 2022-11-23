#!/bin/bash
clc -s -e eg doc.go tdb_test.go tdb1_test.go tdb2_test.go tdb3_test.go
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st
