#!/bin/sh
export GOOS=linux
export GOARCH=amd64
go build main.go
unset GOOS
unset GOARCH
