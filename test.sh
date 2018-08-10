#!/bin/bash
set -e

verbose=""

if [[ "$@" == "-v" ]]
then
  verbose="-v"
fi

# export GO111MODULE=on
# export GO_BIN=go111
go test -tags sqlite -vet off $verbose ./...
