<p align="center"><img src="https://github.com/gobuffalo/buffalo/blob/master/logo.svg" width="360"></p>

<p align="center">
<a href="https://godoc.org/github.com/gobuffalo/buffalo"><img src="https://godoc.org/github.com/gobuffalo/buffalo?status.svg" alt="GoDoc" /></a>
<a href="https://dev.azure.com/markbates/buffalo/_build?definitionId=1"><img src="https://dev.azure.com/markbates/buffalo/_apis/build/status/gobuffalo.buffalo?branchName=master" alt="Build Status" /></a>
<a href="https://goreportcard.com/report/github.com/gobuffalo/buffalo"><img src="https://goreportcard.com/badge/github.com/gobuffalo/buffalo" alt="Go Report Card" /></a>
<a href="https://www.codetriage.com/gobuffalo/buffalo"><img src="https://www.codetriage.com/gobuffalo/buffalo/badges/users.svg" alt="Open Source Helpers" /></a>
</p>

# Buffalo

A Go web development eco-system, designed to make your life easier.

Buffalo helps you to generate a web project that already has everything from front-end (JavaScript, SCSS, etc.) to the back-end (database, routing, etc.) already hooked up and ready to run. From there it provides easy APIs to build your web application quickly in Go.

Buffalo **isn't just a framework**; it's a holistic web development environment and project structure that **lets developers get straight to the business** of, well, building their business.

> I :heart: web dev in go again - Brian Ketelsen

## Documentation

Please visit [http://gobuffalo.io](http://gobuffalo.io) for the latest documentation, examples, and more.

### Quick Start
* [Installation](http://gobuffalo.io/docs/installation)
* [Create a new project](http://gobuffalo.io/docs/new-project)
* [Examples](http://gobuffalo.io/docs/examples)

## Shoulders of Giants

Buffalo would not be possible if not for all of the great projects it depends on. Please see [SHOULDERS.md](SHOULDERS.md) to see a list of them.

### Templating

[github.com/gobuffalo/plush](https://github.com/gobuffalo/plush) - This templating package was chosen over the standard Go `html/template` package for a variety of reasons. The biggest of which is that it is significantly more flexible and easy to work with.

### Routing

[github.com/gorilla/mux](https://github.com/gorilla/mux) - This router was chosen because of its stability and flexibility. There might be faster routers out there, but this one is definitely the most powerful!

### Task Runner (Optional)

[github.com/markbates/grift](https://github.com/markbates/grift) - If you're familiar with Rake tasks from Ruby, you'll be right at home using Grift. This package was chosen to allow for the easy running of simple, and common, tasks that most web applications need. Think things like seeding a database or taking in a CSV file and generating database records. Buffalo ships with an example `routes` task that prints of the defined routes and the function that handles those requests.

### Models/ORM (Optional)

[github.com/gobuffalo/pop](https://github.com/gobuffalo/pop) - Accessing databases is nothing new in web applications. Pop, and its command line tool, Soda, were chosen because they strike a nice balance between simplifying common tasks, being idiomatic, and giving you the flexibility you need to build your app. Pop and Soda share the same core philosophies as Buffalo, so they were a natural choice.

### Sessions, Cookies, WebSockets, and more...

[github.com/gorilla](https://github.com/gorilla) - The Gorilla toolkit is a great set of packages designed to improve upon the standard library for a variety of web-related packages. With these high-quality packages Buffalo can keep its "core" code to a minimum and focus on its goal of gluing them all together to make your life better.

## Benchmarks

Oh, yeah, everyone wants benchmarks! What would a web framework be without its benchmarks? Well, guess what? I'm not giving you any! That's right. This is Go! I assure you that it is plenty fast enough for you. If you want benchmarks you can either a) check out any benchmarks that the [GIANTS](SHOULDERS.md) Buffalo is built upon having published, or b) run your own. I have no interest in playing the benchmark game, and neither should you.

## Contributing

First, thank you so much for wanting to contribute! It means so much that you care enough to want to contribute. We appreciate every PR from the smallest of typos to the be biggest of features.

To contribute, please read the contribution guidelines: [CONTRIBUTING](.github/CONTRIBUTING.md)
