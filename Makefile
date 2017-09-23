#!/usr/bin/make -f

VERSION=$(shell git describe --tags --always)
IMAGE=nickschuch/watson

release: build push

build:
	docker build -t ${IMAGE}:${VERSION} .

push:
	docker push ${IMAGE}:${VERSION}
