#!/usr/bin/env bash

go test ./... -mod=vendor -short -cover -count=1 -race -failfast -timeout 30s
