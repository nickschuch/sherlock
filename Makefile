#!/usr/bin/make -f

VERSION=$(shell git describe --tags --always)
IMAGE=nickschuch/sherlock

release: build push

build:
	docker build -t ${IMAGE}:${VERSION} .

push:
	docker push ${IMAGE}:${VERSION}
