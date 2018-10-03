FROM gobuffalo/buffalo:development

ARG CODECOV_TOKEN
ARG CI
ARG TRAVIS
ARG TRAVIS_BRANCH
ARG TRAVIS_COMMIT
ARG TRAVIS_JOB_ID
ARG TRAVIS_JOB_NUMBER
ARG TRAVIS_OS_NAME
ARG TRAVIS_PULL_REQUEST
ARG TRAVIS_PULL_REQUEST_SHA
ARG TRAVIS_REPO_SLUG
ARG TRAVIS_TAG

ENV BP=$GOPATH/src/github.com/gobuffalo/buffalo

RUN rm $(which buffalo)
RUN rm -rf $BP
RUN mkdir -p $BP
WORKDIR $BP
COPY . .

RUN make ci-deps

RUN packr clean
RUN gometalinter --vendor --deadline=5m ./... --skip=internal
RUN make install

RUN buffalo version

RUN go test -tags "sqlite integration_test" -race  ./...
RUN go test -tags "sqlite integration_test" -coverprofile cover.out -covermode count ./...

RUN if [ -z "$CODECOV_TOKEN"  ] ; then \
    echo codecov not enabled ; \
    else curl -s https://codecov.io/bash -o codecov && \
    bash codecov -f cover.out -X fix; fi

WORKDIR $GOPATH/src/

# START: tests bins are built with tags properly
RUN mkdir -p $GOPATH/src/github.com/markbates
WORKDIR $GOPATH/src/github.com/markbates
RUN buffalo new --skip-webpack coke --db-type=sqlite3
WORKDIR $GOPATH/src/github.com/markbates/coke
RUN buffalo db create -a -d
RUN buffalo g resource widget name
RUN buffalo b -d
# works fine:
RUN ./bin/coke migrate
RUN rm -rfv $GOPATH/src/github.com/markbates/coke
# :END

WORKDIR $GOPATH/src/

RUN buffalo new  --db-type=sqlite3 hello_world --ci-provider=travis
WORKDIR ./hello_world

RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/new_travis.json

RUN go vet ./...
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
