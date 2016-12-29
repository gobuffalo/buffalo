# Buffalo [![GoDoc](https://godoc.org/github.com/gobuffalo/buffalo?status.svg)](https://godoc.org/github.com/gobuffalo/buffalo) [![Build Status](https://travis-ci.org/gobuffalo/buffalo.svg?branch=master)](https://travis-ci.org/gobuffalo/buffalo) [![Code Climate](https://codeclimate.com/github/gobuffalo/buffalo/badges/gpa.svg)](https://codeclimate.com/github/gobuffalo/buffalo)

Buffalo is Go web framework. Yeah, I hate the word "framework" too! Buffalo is different though. Buffalo doesn't want to re-invent wheels like routing and templating. Buffalo is glue that wraps all of the best packages available and makes them all play nicely together.

Buffalo is "idiomatic", for whatever that is worth. The purpose of a framework is not to bend you to it's will, but to rather get out of your way and make your job of building your application easy. That is the goal of Buffalo.

If you were to look through the Buffalo code base you'll find little code, just enough to assemble the amazing packages that other's have written into one coherent system.

> I :heart: web dev in go again - Brian Ketelsen

## Documentation

Please visit [http://gobuffalo.io](http://gobuffalo.io) for the latest documentation, examples, and more.

## Installation

```text
$ go get -u github.com/gobuffalo/buffalo/buffalo
```

## Generating a new Project

Buffalo aims to make building new web applications in Go as simple as possible, and what could be more simple that a new application generator?

```text
$ buffalo new <name>
```

That will generate a whole new Buffalo application that is ready to go. It'll even run `go get` for you to make sure you have all of the necessary dependencies needed to run your application.

To see a list of available flags for the `new` command, just check out it's help.

```text
$ buffalo help new
```

## Running your application

Buffalo is Go "standards" compliant, that means you can just build your binary and run it. It's that simple.

### Running your application in Development

One of the downsides to Go development is the lack of code "reloading". This means as you change your code you need to manually stop your application, rebuild it, and then restart it. Buffalo finds this is annoying, and wants to make life better for you.

```text
$ buffalo dev
```

The `dev` command will watch your `.go` and `.html` files by default and rebuild, and restart, your binary for you so you don't have to worry about such things. Just run the `dev` command and start coding.

## Testing your application

Just like running your application, Buffalo doesn't stop you from using the standard Go tools for testing. Buffalo does ship with a `test` command that will run all of your tests will convienently skipping that pesky old `./vendor` directory!

```text
$ buffalo test
```

## Shoulders of Giants

Buffalo would not be possible if not for all of the great projects it depends on. Please see [SHOULDERS.md](SHOULDERS.md) to see a list of them.

### Templating

[github.com/aymerick/raymond](https://github.com/aymerick/raymond) - This templating package was chosen over the standard Go `html/template` package for a variety of reasons. The biggest of which is that it is significantly more flexible and easier to work with. It also has the added factor of being familiar to those who have worked with "Handlebars" or "Mustache" templates before.

### Routing

[github.com/gorilla/mux](https://github.com/gorilla/mux) - This router was chosen because of it's stability and flexibility. There might be faster routers out there, but this one is definitely the most powerful!

### Task Runner (Optional)

[github.com/markbates/grift](https://github.com/markbates/grift) - If you're familiar with Rake tasks from Ruby, you'll be right at home using Grift. This package was chosen to allow for the easy running of simple, and common, tasks that most web applications need. Think things like seeding a database, or taking in a CSV file and generating database records. Buffalo shops with an example `routes` task that prints off the defined routes and the function that handles those requests.

### Models/ORM (Optional)

[github.com/markbates/pop](https://github.com/markbates/pop) - Accessing databases is nothing new in web applications. Pop, and it's command line tool, Soda, were chosen because they strike a nice balance between simplifying common tasks, being idiomatic, and giving you the flexibility you need to built your app. Pop, and Soda, share the same core philosphies as Buffalo, so they were a natural choice.

### Sessions, Cookies, Websockets, and more...

[github.com/gorilla](https://github.com/gorilla) - The Gorilla toolkit is a great set of packages designed to improve upon the standard libary for a variety of web related packages. With these high quality packages Buffalo is able to keep its "core" code to a minimum and focus on it's goal of glueing them all together to make your life better.

## Benchmarks

Oh, yeah, everyone wants benchmarks! What would a web framework be without it's benchmarks? Well, guess what? I'm not giving you any! That's right. This is Go! I assure you that it is plenty fast enough for you. If you want benchmarks you can either a) checkout any benchmarks that the [GIANTS](SHOULDERS.md) Buffalo is built upon have published, or b) run your own. I have no interest in playing the benchmark game, and neither should you.
