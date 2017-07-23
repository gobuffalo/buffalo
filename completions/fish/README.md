# Install

## Manually

1. See where your `fish` configuration stores completion scripts

```
~ echo $fish_complete_path
/Users/hasitmistry/.config/fish/completions /usr/local/etc/fish/completions /usr/local/share/fish/completions /Users/hasitmistry/.local/share/fish/generated_completions
```

2. Copy `buffalo.fish` file to one of those directories; preferably into `~/.config/fish/completions`.

```
~ cp $GOPATH/src/github.com/gobuffalo/buffalo/completions/fish/README.md $fish_complete_path[1]
```

## The `grift` way

1. Go to `gobuffalo/buffalo` project directory

```
~ cd $GOPATH/src/github.com/gobuffalo/buffalo/
g/s/g/g/buffalo ╍
```

2. Run `completions:fish` grift

```
~ grift completions:fish
```

Learn more about `grifts` [here](https://github.com/markbates/grift).

# Use

Once you have the `buffalo.fish` file situated in the correct directory, `fish` should be able to tab complete `buffalo` commands. Try typing `buffalo<TAB><TAB>` like so.

```
~ buffalo<TAB><TAB>
build  (Builds a Buffalo binary, including bundling of assets (packr & webpack))
console                                (Runs your Buffalo app in a REPL console)
db                                   (A tasty treat for all your database needs)
destroy                                       (Allows to destroy generated code)
dev                                (Runs your Buffalo app in 'development' mode)
generate                        (A collection of generators to make life easier)
help                                                    (Help about any command)
new                                          (Creates a new Buffalo application)
setup              (Setups a newly created, or recently checked out application)
task                                                     (Runs your grift tasks)
test                                       (Runs the tests for your Buffalo app)
version                                    (Print the version number of buffalo)
```