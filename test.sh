#!/bin/bash
set -e

verbose=""

if [[ "$@" == "-v" ]]
then
  verbose="-v"
fi

go test -tags sqlite $verbose ./...
