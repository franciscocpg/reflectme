SHELL := /bin/bash
COVERAGE_FILE = coverage.out 

up-deps:
	glide up

install-deps:
	glide i

test:
	go test -v

test-coverage:
	go test -v -race -covermode=atomic -coverprofile=${COVERAGE_FILE}

coverage-html: test-coverage
	go tool cover -mode=atomic -html=${COVERAGE_FILE}

send-codecov-ci:
	bash <(curl -s https://codecov.io/bash)

send-codecov-local:
	bash <(curl -s https://codecov.io/bash) -t ${CODECOV_TOKEN}
