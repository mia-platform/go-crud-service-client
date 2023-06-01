# Golang CRUD service client

## Introduction

In our projects, we often use [crud-service](https://github.com/mia-platform/crud-service)
and we want a way to interact to it with a standard client.

At the moment, it is limited.
The supported methods are:

- Export (`/export` endpoint)

If you need some other method, please add it with a PR.

## Usage

To use it, install with

```sh
go get github.com/mia-platform/crud-service-client
```

## Development

To run tests:

```sh
make test
```

To generate coverage report:

```sh
make coverage
```
