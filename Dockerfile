FROM gobuffalo/buffalo:development

RUN buffalo version

RUN go get -v -u github.com/golang/lint/golint
RUN go get -v -u github.com/markbates/filetest
RUN go get -v -u github.com/gobuffalo/makr
RUN go get -v -u github.com/markbates/grift
RUN go get -v -u github.com/markbates/inflect
RUN go get -v -u github.com/markbates/refresh
RUN go get -v -u github.com/gobuffalo/tags
RUN go get -v -u github.com/gobuffalo/pop
RUN go get -v -u github.com/mattn/go-sqlite3

ENV BP=$GOPATH/src/github.com/gobuffalo/buffalo

RUN rm $(which buffalo)
RUN rm -rf $BP
RUN mkdir -p $BP
WORKDIR $BP
ADD . .

RUN go get -v -t ./...

RUN go install -v -tags sqlite ./buffalo

RUN go test -tags sqlite -race $(go list ./... | grep -v /vendor/)

RUN golint -set_exit_status $(go list ./... | grep -v /vendor/)


WORKDIR $GOPATH/src/
RUN buffalo new  --db-type=sqlite3 hello_world --ci-provider=travis
WORKDIR ./hello_world

RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/new_travis.json

RUN go vet -x $(go list ./... | grep -v /vendor/)
RUN buffalo db create -a
RUN buffalo db migrate -e test
RUN buffalo test -race

RUN go get -v github.com/gobuffalo/buffalo-goth
RUN buffalo g goth facebook twitter linkedin github
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/goth.json

RUN buffalo g resource admins --skip-model
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_skip_model.json
RUN rm actions/admins_test.go

RUN buffalo test -race
RUN buffalo build -static

RUN buffalo g resource users name:text email:text
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_model_migration.json

RUN buffalo g resource admins --use-model users
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_use_model.json

RUN rm actions/admins_test.go
RUN rm models/user_test.go
RUN rm models/user.go
RUN rm actions/users_test.go
RUN rm -rv templates/users

RUN buffalo g resource ouch
RUN buffalo d resource -y ouch
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_resource_all.json

RUN buffalo db g model ouch
RUN buffalo db d model -y ouch
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_model_all.json

RUN buffalo g actions ouch build edit
RUN buffalo d action -y ouch
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/destroy_action_all.json

RUN buffalo g actions comments show edit
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_all.json

RUN buffalo g actions comments destroy
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_existing.json

RUN buffalo g resource user
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_singular.json

RUN buffalo g resource cars
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_plural.json

RUN buffalo g actions users create --skip-template
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_skip_template.json

RUN buffalo g actions users update --skip-template --method POST
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_action_with_method.json

WORKDIR $GOPATH/src
RUN buffalo new  --api apiapp
WORKDIR ./apiapp
RUN buffalo build
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/apiapp.json

RUN buffalo g task plainTask
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_plain_task.json

RUN buffalo g task nested:task:now
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_nested_task.json

RUN buffalo g resource admin/planes
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_nested.json

RUN buffalo g resource admin/users --name=AdminUser
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_nested_model_name.json

RUN buffalo g resource person
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_resource_irregular.json

RUN buffalo g resource person_event
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_underscore.json

RUN buffalo g mailer welcome_email
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/generate_mailer.json

RUN rm -rf bin
RUN buffalo build -k -e
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/no_assets_build.json
