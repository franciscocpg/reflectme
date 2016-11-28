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

next-version: git-semver
	git semver next

release: next-version
	git push --tags

git-semver:
	git semver 1>/dev/null 2>&1 || (git clone https://github.com/markchalloner/git-semver.git /tmp/git-semver && cd /tmp/git-semver && git checkout $( \
    git tag | grep '^[0-9]\+\.[0-9]\+\.[0-9]\+$' | \
    sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -n 1 \
) && ./install.sh)