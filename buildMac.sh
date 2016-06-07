#!/bin/sh
export GOOS=darwin
export GOARCH=amd64
go build main.go
unset GOOS
unset GOARCH
