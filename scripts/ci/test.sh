#!/bin/bash

set -euo pipefail
set +x
[[ ${DEBUG-1} ]] && set -x

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
source $HERE/defaults.sh

envfile=$HERE/test-env

if [[ ${BUILDKITE_MESSAGE:-} =~ "deploy" ]]
then
  echo "--- Skipping tests"
else
  echo "--- Run tests"
  docker run --rm --env-file $envfile $DOCKER_IMAGE make test
fi
