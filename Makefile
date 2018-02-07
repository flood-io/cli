GIT_SHA = $(shell git rev-parse HEAD)
GITHUB_TOKEN = $(shell cat ~/.github-token)

test-focus : TEST_PACKAGE = config
test-focus : TEST_ARGS = #-run client
test-focus :
	go test github.com/flood-io/cli/${TEST_PACKAGE} -v -p 1 ${TEST_ARGS}

test :
	go test ./...

release:
	goreleaser

build-docker :
	@docker build --build-arg GIT_SHA=${GIT_SHA} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t cli .

ci-local: 
	export BUILDKITE_COMMIT=${GIT_SHA} ; \
	export GITHUB_TOKEN=${GITHUB_TOKEN} ; \
	export NO_PUSH=1 ; \
	export DEBUG=1 ; \
	./scripts/ci/build.sh &&\
	./scripts/ci/test.sh &&\
	./scripts/ci/release.sh
