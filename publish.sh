#!/usr/bin/env bash

source version.txt

jfrog rt go build go --build-name=go-commons --build-number=$IMAGE_VERSION
jfrog rt go-publish go_local $IMAGE_VERSION --build-name=go-commons --build-number=$IMAGE_VERSION --deps=ALL
jfrog rt build-collect-env go-commons $IMAGE_VERSION
jfrog rt build-publish go-commons $IMAGE_VERSION