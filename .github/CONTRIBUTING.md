# Contributing to Buffalo

First, thank you so much for wanting to contribute! It means so much that you care enough to want to contribute. We appreciate every PR from the smallest of typos to the be biggest of features.

## First time committing to a Go Repo?

Contributing to a Go project for the first time can be confusing due to import statements requiring a very specific path on disk. You can use the following two approaches to contribute.

### 1. Using Go Path
If you are using your Go Path for development, please take a moment and read this fantastic post on how to easily work with Go repos.

[https://splice.com/blog/contributing-open-source-git-repositories-go/](https://splice.com/blog/contributing-open-source-git-repositories-go/)

### 2. Using Go Modules
Support for Go Modules were introduced in [Release v0.13.0-beta.1](https://github.com/gobuffalo/buffalo/releases/tag/v0.13.0-beta.1), and you can now use them to easily set up a development environment. The steps involve:

1. Fork the repo
2. Clone the repo to any location in your work station
3. Add a `GO111MODULE` environment variable with `export GO111MODULE=on`
4. Run `make deps` to install dependencies
5. Read the contribution guideline below

## Contributing to Buffalo

1. Check [https://github.com/gobuffalo/buffalo/issues](https://github.com/gobuffalo/buffalo/issues) to make sure you're not working on a duplicate issue or PR.
2. If you want to implement a new feature that doesn't have an issue open, please open one and ask for feedback on the feature before spending a lot of time working on it. It's possible the feature has already been discussed, or it's out of scope, or some other reason that might later prevent a PR from being accepted. The [#buffalo](https://gobuffalo.io/docs/slack) channel on gophers.slack.com is a great place to seek this kind of guidance.
3. Write your feature/fix and make sure to include tests. Tests are an **absolute** requirement for any pull request. Please make sure to use the same testing style and libraries as the rest of the tests.
4. Make sure tests run when doing `go test ./...`. You may need to do `go get -t ./...` first to get the testing dependencies.
5. (Optional) There is a much longer set of integration tests that can be run. These will be run by Travis-CI when you open a PR. If you want to run them locally, you can by running `docker build .`.

Feel free to ask for help, but don't target a specific person (unless you're replying to this person). e.g. don't @ markbates, but @ gobuffalo/core-managers instead.

### Making your Pull Request

Open a PR against the `development` branch. **_DO NOT_** open one against `master`. All "unreleased" work happens in the `development` branch, and we will fix the master branch if necessary.

**WE WILL CLOSE ANY PR OPENED AGAINST MASTER BRANCH**.

## Documentation Welcome

Hands down the most important and the most welcome pull requests are for documentation. We LOVE documentation PRs, and so do all those that come after you.

Whether it's GoDoc or prose on [http://gobuffalo.io](http://gobuffalo.io) all documentation is welcome.

You can submit PRs to change the website and/or docs on [https://github.com/gobuffalo/gobuffalo](https://github.com/gobuffalo/gobuffalo).

## Thank You

Once again, we want to take the chance to say thank you again for wanting to contribute to Buffalo. This is a community project, and that means we **need** your help! Thank you so much.
