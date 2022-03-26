#!/usr/bin/env bash

go test ./... -mod=vendor -short -cover -count=1 -v -race -failfast -timeout 300s
