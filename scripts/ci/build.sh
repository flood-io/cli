#!/bin/bash

set -euo pipefail
set +x
[[ ${DEBUG-1} ]] && set -x

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
source $HERE/defaults.sh

echo "--- Docker build"
docker build --build-arg GIT_SHA=${SHORT_SHA}  --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t $DOCKER_IMAGE .

echo "--- Docker push"
if [[ $LOCAL_ONLY ]]; then
  echo "(not pushing)"
else
  docker push $DOCKER_IMAGE
fi

echo "--- Check test dependencies"
docker run --rm $DOCKER_IMAGE goreleaser help  1> /dev/null && echo OK || exit 1
