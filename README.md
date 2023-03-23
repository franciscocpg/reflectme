ReflectME
===========
![example workflow](https://github.com/franciscocpg/reflectme/actions/workflows/main.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/franciscocpg/reflectme.svg)](https://pkg.go.dev/github.com/franciscocpg/reflectme)

Package reflectme provides high level abstractions above the golang reflect library with some utilities functions.

The base code is a fork of https://github.com/oleiade/reflections. Then some other concepts and functions were added.


## Installation
```shell
go get github.com/franciscocpg/reflectme@latest
```

## Usage
Take a look at [tests](https://github.com/franciscocpg/reflectme/blob/master/reflections_test.go) and [go reference](https://pkg.go.dev/github.com/franciscocpg/reflectme).

## Contribute

* Check for open issues or open a fresh issue to start a discussion around a feature idea or a bug.
* Fork the repository on GitHub to start making your changes to the **master** branch (or branch off of it).
* Write tests which shows that the bug was fixed or that the feature works as expected and make sure the is 100% coverage. One can run `make test-coverage` to see the coverage %. If it's not 100%, one can run `make coverage-missing` to catch the lines that are not covered.
* Send a pull request
* If all checks status are successful the PR is going to be merged.
