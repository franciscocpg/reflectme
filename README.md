ReflectME
===========
[![Build Status](https://travis-ci.org/franciscocpg/reflectme.svg?branch=master)](https://travis-ci.org/franciscocpg/reflectme)
[![GoDoc](https://godoc.org/github.com/franciscocpg/reflectme?status.svg)](https://godoc.org/github.com/franciscocpg/reflectme)

Package reflectme provides high level abstractions above the golang reflect library with some utilities functions.

The base code is a fork of https://github.com/oleiade/reflections. Then some other concepts and functions were added.


## Installation
Use some dependency manager tool like [glide](https://github.com/Masterminds/glide) pinning to some tag version. Eg
```yaml
package: github.com/franciscocpg/test
import:
- package: github.com/franciscocpg/reflectme
  version: 0.1.4
```

## Usage
As the code is 100% test covered the easier way is to look at [tests](https://github.com/franciscocpg/reflectme/blob/master/reflections_test.go) and [godoc](https://godoc.org/github.com/franciscocpg/reflectme).

## Contribute

* Check for open issues or open a fresh issue to start a discussion around a feature idea or a bug.
* Fork the repository on GitHub to start making your changes to the **master** branch (or branch off of it).
* Write tests which shows that the bug was fixed or that the feature works as expected and make sure the is 100% coverage. One can run `make test-coverage` to see the coverage %. If it's not 100%, one can run `make coverage-missing` to catch the lines that are not covered.
* Send a pull request
* If all checks status are successful the PR is going to be merged.
* Every merged PR is going to tag a new `semver` release