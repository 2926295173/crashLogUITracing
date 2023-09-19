#!/bin/sh
appname="clash-log-record"

author="openrhc-"

version=$author`date +"%Y%m%d"`

CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o $appname -trimpath main.go

# scp clash-log-record root@192.168.0.2:/root