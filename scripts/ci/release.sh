#!/bin/bash

set -euo pipefail
set +x
[[ ${DEBUG-1} ]] && set -x

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
source $HERE/defaults.sh

if [[ $BUILDKITE_TAG ]]; then
  if [[ -z $LOCAL_ONLY ]]; then
    docker pull $DOCKER_IMAGE
  fi
  docker run --rm $DOCKER_IMAGE make release
fi
