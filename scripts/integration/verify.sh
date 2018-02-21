#!/bin/bash

set -euo pipefail

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
cd $HERE/../..


: ${BUGSNAG_API_KEY:=}
export FLOOD_DEBUG=1
go run -ldflags "-X main.bugsnagAPIKey=${BUGSNAG_API_KEY}" *.go verify --host http://localhost:5000 --verbose "$@"

