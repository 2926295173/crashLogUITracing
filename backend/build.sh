#!/bin/sh
appname="clash-tracing.exe"

author="openrhc-"

version=$author`date +"%Y%m%d"`

go build -o $appname -trimpath -ldflags "-s -w -buildid=" main.go