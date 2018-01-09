.PHONY: proto proto-deps

proto-deps:
	go get github.com/gogo/protobuf
	go get github.com/gogo/protobuf/proto
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/gogo/protobuf/gogoproto

proto:
	cd proto ; \
	protoc -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf --gofast_out=plugins=grpc:. *.proto
