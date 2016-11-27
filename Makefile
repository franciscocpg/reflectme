SHELL := /bin/bash

up-deps:
	glide up

install-deps:
	glide i

test:
	go test -v
