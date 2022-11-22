#!/bin/bash
clc -s -e doc.go csv_test.go classic_test.go tdb_test.go db1_test.go eg bin
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st
