#!/bin/bash

set -ex

# git branch --set-upstream-to=origin/$TRAVIS_BRANCH $TRAVIS_BRANCH
BP=$GOPATH/src/github.com/gobuffalo/buffalo

export GO111MODULE=on

go get github.com/markbates/filetest

make install
buffalo version
go test -tags "sqlite integration_test" -race  ./...

cd $GOPATH/src/

# START: tests bins are built with tags properly
mkdir -p $GOPATH/src/github.com/markbates
cd $GOPATH/src/github.com/markbates
buffalo new --skip-webpack coke --db-type=sqlite3
cd $GOPATH/src/github.com/markbates/coke
buffalo db create -a -d
buffalo g resource widget name
buffalo b
# works fine:
./bin/coke migrate
rm -rfv $GOPATH/src/github.com/markbates/coke
# :END

cd $GOPATH/src/

buffalo new --db-type=sqlite3 hello_world --ci-provider=travis
cd ./hello_world

filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/new_travis.json

go vet ./...
buffalo db create -a
buffalo db migrate -e test
buffalo test -race

buffalo g resource admins --skip-model
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_skip_model.json
rm actions/admins_test.go

buffalo test -race
buffalo build -static

buffalo g resource users name:text email:text
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_model_migration.json

buffalo g resource admins --use-model users
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_use_model.json

rm actions/admins_test.go
rm models/user_test.go
rm models/user.go
rm actions/users_test.go
rm -rv templates/users

buffalo g resource ouch
buffalo d resource -y ouch
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_resource_all.json

buffalo db g model ouch
buffalo db d model -y ouch
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_model_all.json

buffalo g actions ouch build edit
buffalo d action -y ouch
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_action_all.json

buffalo g mailer ouch
buffalo d mailer -y ouch
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_mailer_all.json

buffalo g actions comments show edit
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_all.json

buffalo g actions comments destroy
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_existing.json

buffalo g resource user
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_singular.json

buffalo g resource cars
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_plural.json

buffalo g resource admin/planes
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_nested_web.json

buffalo g resource admin/users --name=AdminUser
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_nested_model_name_web.json

buffalo g actions users create --skip-template
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_skip_template.json

buffalo g actions users update --skip-template --method POST
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_with_method.json

cd $GOPATH/src
buffalo new --api apiapp
cd ./apiapp
buffalo build
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/apiapp.json

buffalo g task plainTask
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_plain_task.json

buffalo g task nested:task:now
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_nested_task.json

buffalo g resource admin/planes
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_nested_api.json

buffalo g resource admin/users --name=AdminUser
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_nested_model_name_api.json

buffalo g resource person
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_irregular.json

buffalo g resource person_event
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_underscore.json

buffalo g mailer welcome_email
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_mailer.json

rm -rf bin
buffalo build -k -e
filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/no_assets_build.json

