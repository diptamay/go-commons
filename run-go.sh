#!/usr/bin/env bash
# This is clone from /usr/local/cal/bin/node script for go with some modification
# you can update env vars as needed
# Usage: sh ./run-go.sh go test ./... -short -v -cover -failfast -timeout 300s

set -e
set -o pipefail

GO_IMAGE=<docker_url>/core/golang-buildtools:1.13.8

INTERACTIVE_FLAGS=''
if [ -t 0 ]; then
  INTERACTIVE_FLAGS='-it'
fi
ENDPOINT_VARS=($(env | grep "^[A-Z]*_ENDPOINT="))

echo -e "Using image '${GO_IMAGE}'\n"
docker pull $GO_IMAGE
docker run \
  --rm \
  $INTERACTIVE_FLAGS \
  -e "CGO_ENABLED=1" \
  -e "GOPATH=/srv/package" \
  -e "LOG_LEVEL=debug" \
  -e "AWS_REGION=us-east-1" \
  -e "IGNORE_EXECUTION_ERROR_NOT_FOUND=true" \
  -v $PWD:/srv/package/src/github.com/diptamay/go-commons \
  -w /srv/package/src/github.com/diptamay/go-commons  \
  $(printf -- " -e %s" ${ENDPOINT_VARS[*]}) \
  $GO_IMAGE "$@"