COVERAGE_FILE = coverage.out
export GOBIN := $(PWD)/bin

gocov := $(GOBIN)/gocov
$(gocov):
	@go install github.com/axw/gocov/gocov@latest

test:
	go test -v

test-coverage:
	go test -v -race -covermode=atomic -coverprofile=${COVERAGE_FILE}

coverage-html: test-coverage
	go tool cover -html=${COVERAGE_FILE}

coverage-missing: $(gocov) test-coverage
	$(gocov) convert ${COVERAGE_FILE} | $(gocov) annotate - | grep MISS

send-codecov-local:
	bash <(curl -s https://codecov.io/bash) -t ${CODECOV_TOKEN}
