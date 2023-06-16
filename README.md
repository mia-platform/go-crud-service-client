[![Test and build](https://github.com/mia-platform/go-crud-service-client/actions/workflows/test-builds.yml/badge.svg)](https://github.com/mia-platform/go-crud-service-client/actions/workflows/test-builds.yml)
[![Coverage Status](https://coveralls.io/repos/github/mia-platform/go-crud-service-client/badge.svg?branch=main)](https://coveralls.io/github/mia-platform/go-crud-service-client?branch=main)
[![GoDoc](https://godoc.org/github.com/davidebianchi/gswagger?status.svg)](https://godoc.org/github.com/davidebianchi/gswagger)

# Golang CRUD service client

## Introduction

In our projects, we often use [CRUD Service](https://github.com/mia-platform/crud-service)
and we want a way to interact to it with a standard client.

At the moment, it is limited.
The supported methods are:

- GetById: `GET`
- List: `GET /`
- Count: `GET /count`
- Export: `GET /export`
- PatchById: `PATCH /:id`
- Patch: `PATCH /`
- Create: `POST /`
- DeleteById: `DELETE /:id`
- Delete: `DELETE /`

If you need some other method, please add it with a PR.

## Usage

To use it, install with

```sh
go get github.com/mia-platform/go-crud-service-client
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
