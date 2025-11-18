#!/bin/bash
export GOOS=linux
export CGO_ENABLED=1
version=`date "+%y.%m.%d %H:%M"`
echo Setup Version: $version
echo package core > ../../core/version.go
echo const VERSION = \"${version}\" >> ../../core/version.go
go build -o registry -ldflags="-s -w"