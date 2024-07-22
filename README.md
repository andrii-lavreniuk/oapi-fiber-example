# oapi-fiber-example

REST API implementation with code generation from OpenAPI specification and Fiber framework

### Before starting install development tools used to build, test and run application

```bash
$ make install
```

The following tools will be installed:

  - [oapi-codegen](https://github.com/deepmap/oapi-codegen) is a command-line tool and library to convert OpenAPI specifications to Go code, be it server-side implementations, API clients, or simply HTTP models
  - [wire](https://github.com/google/wire) is a code generation tool that automates connecting components using dependency injection
  - [mockery](https://vektra.github.io/mockery/latest/) is a project that creates mock implementations of Golang interfaces
  - [air](https://github.com/air-verse/air) is a live-reloading command line utility for developing Go applications
  - [golangci-lint](https://github.com/golangci/golangci-lint) is a fast Go linters runner

### Makefile

Check the Makefile opportunities and available commands:

```bash
$ make help
```

### Generate code

```bash
# generate server code from openapi specification
$ make api
# generate DI with wire
$ make wire
# generate mocks
$ make mocks
```

### Run service

Copy and setup `.env` file:
```bash
$ cp .env.example .env
```

Run database:
```bash
$ make compose
```

Run migrations to create table and insert test data:
```bash
$ go run cmd/cli/main.go db init
$ go run cmd/cli/main.go db migrate

# see all available commands
$ go run cmd/cli/main.go db
```

Run with [air](https://github.com/cosmtrek/air):

```bash
$ make run
```

Check the API:
```bash
# get profiles list
$ curl --location 'localhost:8080/v1/profiles' --header 'X-Api-Key: www-dfq92-sqfwf'

# get profile by username
$ curl --location 'localhost:8080/v1/profiles?username=guest' --header 'X-Api-Key: www-dfq92-sqfwf'
```

## Run linter

```bash
$ make lint
```

### Run tests

```bash
$ make test
```