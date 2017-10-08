#!/bin/bash

set -ex

go get -v -t $(go list ./... | grep -v /vendor/)

go install -v ./buffalo

go test -race $(go list ./... | grep -v /vendor/)

golint -set_exit_status $(go list ./... | grep -v /vendor/)


cd $GOPATH/src/
buffalo new  --db-type=sqlite3 hello_world --ci-provider=travis
cd ./hello_world

filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/new_travis.json

go vet -x $(go list ./... | grep -v /vendor/)
buffalo db create -a
buffalo db migrate -e test
buffalo test -race

go get -v github.com/gobuffalo/buffalo-goth
buffalo g goth facebook twitter linkedin github
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/goth.json

buffalo g resource admins --skip-model
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_skip_model.json
rm actions/admins_test.go

buffalo test -race
buffalo build -static

buffalo g resource users name:text email:text
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_model_migration.json

rm models/user_test.go
rm models/user.go
rm actions/users_test.go
rm -rv templates/users

buffalo g resource --type=json users name:text email:text
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_json-xml.json

rm models/user_test.go
rm models/user.go
rm actions/users_test.go

buffalo g resource --type=xml users name:text email:text
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_json-xml.json

rm models/user_test.go
rm models/user.go
rm actions/users_test.go

buffalo g resource ouch
buffalo d resource -y ouch
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_resource_all.json

buffalo db g model ouch
buffalo db d model -y ouch
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_model_all.json

buffalo g actions ouch build edit
buffalo d action -y ouch
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_action_all.json

buffalo g actions comments show edit
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_all.json

buffalo g actions comments destroy
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_existing.json

buffalo g resource user
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_singular.json

buffalo g resource cars
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_plural.json

buffalo g actions users create --skip-template
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_skip_template.json

buffalo g actions users update --skip-template --method POST
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_with_method.json

cd $GOPATH/src
buffalo new  --api apiapp

cd ./apiapp
buffalo build
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/apiapp.json

buffalo g task plainTask
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_plain_task.json

buffalo g task nested:task:now
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_nested_task.json

buffalo g resource admin/planes
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_nested.json

buffalo g resource admin/users --model-name=AdminUser
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_nested_model_name.json
