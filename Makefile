SHELL := /bin/bash
COVERAGE_FILE = coverage.out

test:
	go test -v

test-coverage:
	go test -v -race -covermode=atomic -coverprofile=${COVERAGE_FILE}

coverage-html: test-coverage
	go tool cover -html=${COVERAGE_FILE}

coverage-missing: gocov test-coverage
	gocov convert ${COVERAGE_FILE} | gocov annotate - | grep MISS

send-codecov-ci:
	bash <(curl -s https://codecov.io/bash)

send-codecov-local:
	bash <(curl -s https://codecov.io/bash) -t ${CODECOV_TOKEN}

gocov:
	gocov=$(shell which gocov)
	[ ! -z $${gocov} ] || go get -v github.com/axw/gocov/gocov
