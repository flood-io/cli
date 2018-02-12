#!/bin/bash

set -euo pipefail

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
cd $HERE/..

if ! git diff --quiet --ignore-submodules HEAD; then
  echo repo is dirty, cannot continue
  exit 1
fi

bindata=static/init-skeleton/bindata.go
make static

if git diff --quiet --ignore-submodules HEAD; then
  # nothing to do...
  exit 0
fi

git add $bindata
git commit -m'rebuilt static assets'
