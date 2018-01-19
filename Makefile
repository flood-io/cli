test-focus : TEST_PACKAGE = config
test-focus : TEST_ARGS = #-run client
test-focus :
	go test github.com/flood-io/cli/${TEST_PACKAGE} -v -p 1 ${TEST_ARGS}
