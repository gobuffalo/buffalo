#!/bin/bash

set +e

clear

echo "postgres"
SODA_DIALECT=postgres go test -bench=.
echo "--------------------"
echo "mysql"
SODA_DIALECT=mysql go test -bench=.
echo "--------------------"
echo "sqlite"
SODA_DIALECT=sqlite go test -bench=.
