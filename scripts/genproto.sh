#!/bin/bash

set -euo pipefail

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
ROOT=$HERE/..

cd $ROOT

go get github.com/gogo/protobuf/proto
go get github.com/gogo/protobuf/protoc-gen-gofast
go get github.com/gogo/protobuf/gogoproto

cd proto

FLOOD_GH_PREFIX=github.com/flood-io/

if [[ $VIRTUALGO ]]; then
  GOPATH=$VIRTUALGO_PATH
fi

if [[ $GOPATH == *":"* ]]; then
  echo "GOPATH has multiple elements ($GOPATH)"
  echo "and I'm not smart enough to deal"
  exit 1
fi

protoc \
  -I=. \
  -I=${GOPATH}/src \
  -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
  --gofast_out=plugins=grpc,import_prefix=$FLOOD_GH_PREFIX:. \
  *.proto

# pinched from https://github.com/coreos/etcd/blob/master/scripts/genproto.sh
# to get around some protobuf tool shortcomings

# cp control.pb.go control.orig.go

sed -i.bak -E "s/github\.com\/flood-io\/(gogoproto|github\.com|golang\.org|google\.golang\.org)/\1/g" ./*.pb.go
sed -i.bak -E "s/github\.com\/flood-io\/google\/protobuf/github.com\/gogo\/protobuf\/types/g" ./*.pb.go
sed -i.bak -E 's/github\.com\/flood-io\/(errors|fmt|io)/\1/g' ./*.pb.go
sed -i.bak -E 's/import _ \"gogoproto\"//g' ./*.pb.go
sed -i.bak -E 's/import fmt \"fmt\"//g' ./*.pb.go
sed -i.bak -E 's/import _ \"github\.com\/flood-io\/google\/api\"//g' ./*.pb.go
sed -i.bak -E 's/import _ \"google\.golang\.org\/genproto\/googleapis\/api\/annotations\"//g' ./*.pb.go
rm -f ./*.bak
goimports -w ./*.pb.go
