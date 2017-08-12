#!/bin/bash
set -e
clear

verbose=""

echo $@

if [[ "$@" == "-v" ]]
then
  verbose="-v"
fi

docker-compose up -d
sleep 4 # Ensure mysql is online

go build -v -o tsoda ./soda

function test {
  echo "!!!Testing $1"
  export SODA_DIALECT=$1
  echo ./tsoda -v
  ! ./tsoda drop -e $SODA_DIALECT -c ./database.yml
  ! ./tsoda create -e $SODA_DIALECT -c ./database.yml
  ./tsoda migrate -e $SODA_DIALECT -c ./database.yml
  ./tsoda migrate down -e $SODA_DIALECT -c ./database.yml
  ./tsoda migrate down -e $SODA_DIALECT -c ./database.yml
  ./tsoda migrate -e $SODA_DIALECT -c ./database.yml
  go test ./...
}

test "postgres"
test "sqlite"
test "mysql"

docker-compose down

rm tsoda
find . -name *.sqlite* -delete
