# fish completion for buffalo

function __fish_buffalo_grifts
  command grift list $argv ^/dev/null | string match -r '\s.*?#' | string trim -c ' #' | string trim
end

# general
complete -xc 'buffalo' -n '__fish_use_subcommand' -s h -l help -d 'help for buffalo'

# build
complete -xc 'buffalo' -n '__fish_use_subcommand' -a build -d 'Builds a Buffalo binary, including bundling of assets (packr & webpack)'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from build' -s e -l extract-assets -d 'extract the assets and put them in a distinct archive'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from build' -s h -l help -d 'help for build'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from build' -l ldflags -d 'set any ldflags to be passed to the go build'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from build' -s o -l output -d 'set the name of the binary (default "bin/completions")'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from build' -s s -l static -d 'build a static binary using --ldflags \'-linkmode external -extldflags "-static"\' (USE FOR CGO)'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from build' -s t -l tags -d 'compile with specific build tags'

# console
complete -xc 'buffalo' -n '__fish_use_subcommand' -a console -d 'Runs your Buffalo app in a REPL console'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from console' -s h -l help -d 'help for console'

# db
complete -xc 'buffalo' -n '__fish_use_subcommand' -a db -d 'A tasty treat for all your database needs'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db' -s c -l config -d 'The configuration file you would like to use'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db' -s d -l debug -d 'Use debug/verbose mode'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db' -s e -l env -d 'The environment you want to run migrations against. Will use $GO_ENV if set (default "development")'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db' -s h -l help -d 'help for db'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db' -s p -l path -d 'Path to the migrations folder (default "./migrations")'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db' -s v -l version -d 'Show version information'
## db create
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and not __fish_seen_subcommand_from create destroy drop generate migrate schema' -a create -d 'Creates databases for you'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from create' -s a -l all -d 'Creates all of the databases in the database.yml'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from create' -s h -l help -d 'help for create'
## db destroy
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and not __fish_seen_subcommand_from create destroy drop generate migrate schema' -a destroy -d 'Allows to destroy generated code'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from destroy' -s h -l help -d 'help for destroy'
## db drop
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and not __fish_seen_subcommand_from create destroy drop generate migrate schema' -a drop -d 'Drops databases for you'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from drop' -s a -l all -d 'Drops all of the databases in the database.yml'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from drop' -s h -l help -d 'help for drop'
## db generate
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and not __fish_seen_subcommand_from create destroy drop generate migrate schema' -a generate
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from generate' -s h -l help -d 'help for generate'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from create destroy drop fizz model sql migrate schema' -a config -d 'Generates a database.yml file for your project'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from create destroy drop config model sql migrate schema' -a fizz -d 'Generates Up/Down migrations for your database using fizz'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from create destroy drop config fizz sql migrate schema' -a model -d 'Generates a model for your database'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from create destroy drop config fizz model migrate schema' -a sql -d 'Generates Up/Down migrations for your database using SQL'
## db migrate
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and not __fish_seen_subcommand_from create destroy drop generate migrate schema' -a migrate -d 'Runs migration against your database'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from migrate' -s h -l help -d 'help for migrate'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from migrate; and not __fish_seen_subcommand_from create destroy drop generate down reset status schema' -a up -d 'Apply all of the \'up\' migrations.'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from migrate; and not __fish_seen_subcommand_from create destroy drop generate up reset status schema' -a down -d 'Apply one or more of the \'down\' migrations.'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from migrate; and not __fish_seen_subcommand_from create destroy drop generate up down status schema' -a reset -d 'The equivalent of running `migrate down` and then `migrate up`'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from migrate; and not __fish_seen_subcommand_from create destroy drop generate up down reset schema' -a status -d 'Displays the status of all migrations.'
## db schema
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and not __fish_seen_subcommand_from create destroy drop generate migrate schema' -a schema -d 'Tools for working with your database schema'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from schema' -s h -l help -d 'help for schema'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from schema; and not __fish_seen_subcommand_from create destroy drop generate migrate dump' -a load -d 'Load a schema.sql file into a database'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from db; and __fish_seen_subcommand_from schema; and not __fish_seen_subcommand_from create destroy drop generate migrate load' -a dump -d 'Dumps out the schema of the selected database'

# destroy
complete -xc 'buffalo' -n '__fish_use_subcommand' -a destroy -d 'Allows to destroy generated code'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from destroy' -s h -l help -d 'help for destroy'
## desctroy action
complete -xc 'buffalo' -n '__fish_seen_subcommand_from destroy; and not __fish_seen_subcommand_from action resource' -a action -d 'Destroys action files'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from destroy; and __fish_seen_subcommand_from action' -s h -l help -d 'help for action'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from destroy; and __fish_seen_subcommand_from action' -s y -l yes -d 'confirms all beforehand'
## desctroy resource
complete -xc 'buffalo' -n '__fish_seen_subcommand_from destroy; and not __fish_seen_subcommand_from action resource' -a resource -d 'Destroys resource files'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from destroy; and __fish_seen_subcommand_from resource' -s h -l help -d 'help for resource'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from destroy; and __fish_seen_subcommand_from resource' -s y -l yes -d 'confirms all beforehand'

# dev
complete -xc 'buffalo' -n '__fish_use_subcommand' -a dev -d 'Runs your Buffalo app in \'development\' mode'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from dev' -s h -l help -d 'help for dev'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from dev' -s d -l debug -d 'use delve to debug the app'

# generate
complete -xc 'buffalo' -n '__fish_use_subcommand' -a generate -d 'A collection of generators to make life easier'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate' -s h -l help -d 'help for generate'
## generate action
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from action docker goth resource task webpack' -a action -d 'Generates new action(s)'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from action' -s h -l help -d 'help for action'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from action' -s m -l method -d 'allows to set a different method for the action being generated. (default "GET")'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from action' -l skip-template -d 'makes resource generator not to generate template for actions'
## generate docker
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from action docker goth resource task webpack' -a docker -d 'Generates a Dockerfile'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from docker' -s h -l help -d 'help for docker'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from action' -l style -d 'what style Dockerfile to generate [multi, standard] (default "multi")'
## generate goth
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from action docker goth resource task webpack' -a goth -d 'Generates a actions/goth.go file configured to the specified providers.'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from docker' -s h -l help -d 'help for goth'
## generate resource
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from action docker goth resource task webpack' -a resource -d 'Generates a new actions/resource file'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from resource' -s h -l help -d 'help for resource'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from resource' -s n -l model-name -d 'allows to define a different model name for the resource being generated.'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from resource' -s s -l skip-migration -d 'sets resource generator not-to add model migration'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from resource' -l skip-model -d 'makes resource generator not to generate model nor migrations'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from resource' -l type -d 'sets the resource type (html or json) (default "html")'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from resource' -s u -l use-model -d 'generates crud options for a model'
## generate task
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from action docker goth resource task webpack' -a task -d 'Generates a grift task'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from task' -s h -l help -d 'help for task'
## generate webpack
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and not __fish_seen_subcommand_from action docker goth resource task webpack' -a webpack -d 'Generates a webpack asset pipeline.'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from webpack' -s h -l help -d 'help for webpack'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from generate; and __fish_seen_subcommand_from webpack' -l with-yarn -d 'allows the use of yarn instead of npm as dependency manager'

# help
complete -xc 'buffalo' -n '__fish_use_subcommand' -a help -d 'Help about any command'

# new
complete -xc 'buffalo' -n '__fish_use_subcommand' -a new -d 'Creates a new Buffalo application'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -s h -l help -d 'help for new'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -l api -d 'skip all front-end code and configure for an API server'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -l ci-provider -a 'none travis gitlab-ci' -d 'specify the type of ci file you would like buffalo to generate [none, travis, gitlab-ci] (default "none")'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -l db-type -a 'postgres mysql sqlite3' -d 'specify the type of database you want to use [postgres, mysql, sqlite3] (default "postgres")'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -l docker -a 'none multi standard' -d 'specify the type of Docker file to generate [none, multi, standard] (default "multi")'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -s f -l force -d 'delete and remake if the app already exists'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -l skip-dep -d 'skips adding github.com/golang/dep to your app'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -l skip-pop -d 'skips adding pop/soda to your app'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -l skip-webpack -d 'skips adding Webpack to your app'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -s v -l verbose -d 'verbosely print out the go get/install commands'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from new' -l with-yarn -d 'allows the use of yarn instead of npm as dependency manager'

# setup
complete -xc 'buffalo' -n '__fish_use_subcommand' -a setup -d 'Setups a newly created, or recently checked out application'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from setup' -s h -l help -d 'help for setup'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from setup' -s d -l drop -d 'drop existing databases'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from setup' -s u -l update -d 'run go get -u against the application\'s Go dependencies'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from setup' -s v -l verbose -d 'run with verbose output'

# task
complete -xc 'buffalo' -n '__fish_use_subcommand' -a task -d 'Runs your grift tasks'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from task' -s h -l help -d 'help for task'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from task' -a '(__fish_buffalo_grifts)' -d 'Grift'

# test
complete -xc 'buffalo' -n '__fish_use_subcommand' -a test -d 'Runs the tests for your Buffalo app'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from test' -s h -l help -d 'help for test'

# version
complete -xc 'buffalo' -n '__fish_use_subcommand' -a version -d 'Print the version number of buffalo'
complete -xc 'buffalo' -n '__fish_seen_subcommand_from version' -s h -l help -d 'help for version'