# Contributing to Buffalo

First, thank you so much for wanting to contribute! It means so much that you care enough to want to contribute. We appreciate every PR from the smallest of typos to the be biggest of features.

## First time committing to a Go Repo?

Contributing to a Go project for the first time can be confusing due to import statements requiring a very specific path on disk.

Please take moment and read this fantastic post on how to easily work with Go repos.

[https://splice.com/blog/contributing-open-source-git-repositories-go/](https://splice.com/blog/contributing-open-source-git-repositories-go/)

## Contributing to Buffalo

1. Check [https://github.com/gobuffalo/buffalo/issues](https://github.com/gobuffalo/buffalo/issues) to make sure you're not working on a duplicate issue or PR.
2. If you want to implement a new feature that doesn't have an issue open. Please open one and ask for feedback on the feature before spending a lot of time working on it. It's possible the feature has already been discussed, or it's out of scope, or some other reason that might later prevent a PR from being accepted. The #buffalo channel on gophers.slack.com is a great place to seek this kind of guideance.
3. Write your feature/fix and make sure to include tests. Tests are an **absolute** requirement for any pull request. Please make sure to use the same testing style and libraries as the rest of the tests.
4. Make sure tests run when doing `go test ./...`. You may need to do `go get -t ./...` first to get the testing dependencies.5. (Optional) There is a much longer set of integration tests that can be run. These will be run by Travis-CI when you open a PR. If you want to run them locally you can by running `$ docker build .`.
5. Open a PR against the `development` branch. Do not open one against `master` unless you are explicitly told to. All "unreleased" work happens in the `development` branch.

## Documentation Welcome

Hands down the most important, and the most welcome, pull requests are for documentation. We LOVE documentation PRs, and so don't all those that come after you.

Whether it's GoDoc or prose on [http://gobuffalo.io](http://gobuffalo.io) all documentation is welcome.

## Thank You

Once again, we just want to take the chance to say thank you again for wanting to contribute to Buffalo. This is a community project and that means we **need** your help! Thank you so much.
