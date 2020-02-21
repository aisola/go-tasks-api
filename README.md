Tasks API
=========

This API, written in [Go](https://golang.org), exposes a simple API for
managing tasks on a person's todo list. The application is intended to serve as
an example for how an engineer at a large enterprise could structure and build
a microservice-style application using Go.

The structure of the application itself is quite opinionated. In fact,
engineers familiar with developing Java microservices for enterprises may find
it relatively easy to navigate. Much of what is written is based off of Go's
mature standard library. However, there are some additional libraries and tools
which, in my humble opinion, become invaluable to serving enterprise applications.

Some tools this application uses:

- [Chi](http://pkg.go.dev/github.com/go-chi/chi) - An HTTP path Router. The
    HTTP router provided in [Go's standard library](http://pkg.go.dev/net/http)
    is sufficient for naive and simple purposes only.
- [Viper](http://pkg.go.dev/github.com/spf13/viper) - A configuration
    management library allowing you to pull configuration from the environment,
    files, flags, or even remote services.
- [Pflag](http://pkg.go.dev/github.com/spf13/pflag) - A drop-in replacement for
    Go's standard libary flag package which supports the GNU/POSIX command line
    flags that we are all used to. Go's standard library flags package does not
    conform to that standard for _some reason_.
- [SQLx](http://pkg.go.dev/github.com/jmoiron/sqlx) - Extensions to Go's standard
    `database/sql` which makes working with databases a little bit more fun.
- [Sqlite](http://pkg.go.dev/github.com/mattn/sqlite) - For connecting to the
    sqlite database. For a real application, this would obviously be replaced
    by some other more production ready database.
- [Zap](http://pkg.go.dev/go.uber.org/zap) - Blazing fast, zero-allocation
    logging. Written at Uber.
- [Is](http://pkg.go.dev/github.com/matryer/is) - A tiny, but professional test
    assertion library to help remove boilerplate from your tests.

I cannot take full credit for the structure of the application. In fact, it is
fairly heavily inspired by [Mat Ryer](https://medium.com/@matryer/how-i-write-go-http-services-after-seven-years-37c208122831)'s
famous talk about writing HTTP services. The specific library choices are my
own, though most of these are quite commonly used in Go applications in general.

### Requirements

 The application requires having Go 1.13+ installed. Aside from that, this
 should build on any platform supported by Go.

### Running Tests

Running the tests for the application is very easy. Go's toolchain has a test
runner and testing basics in the standard library. In order to run the tests
for the whole application, simply `cd` into the directory that you have cloned
the repository into and run:

```sh
go test ./...
```

### Building and Running the Application

There are two ways to do this based on your needs. The easiest way to get up
and running is to use Go's `run` command. Simply run the following from the
repo root.

```sh
go run cmd/tasks/main.go
```

This will pull all dependencies, build the application, and run the resulting
binary.

The second way is to use Go's `build` command which will build the binary. You
can then run the application by calling the binary at the command line.

```sh
go build ./cmd/tasks
./tasks
```

If you want to change the port that the application is running on, you may
provide the `--bind` and if you'd like to have a persistent sqlite database,
the provide the `--database` flag with a path. You can figure out the details
by running `./tasks --help`.

The API itself is really simple. Reading the code a bit should give you a
decent understanding of what the actual API is. Hint: It's not very
interesting.
