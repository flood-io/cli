
INPUTS := test.ts $(wildcard *.json) yarn.lock
IGNORE := "-ignore=Makefile" "-ignore=.*.log" "-ignore=bindata.go" "-ignore=.gitignore"

bindata.go : ${INPUTS}
	go-bindata ${IGNORE} -pkg initskeleton .

clean:
	rm -f bindata.go

.PHONY: clean
