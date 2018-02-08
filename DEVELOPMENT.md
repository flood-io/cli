# Development

## Note

If you're a curious non-flooder thanks for stopping by, but unfortunately right now you won't be able to build
the `flood` cli yourself.

This is because the code depends on our private `github.com/flood-io/go-wrenches` package.

If non-flood folks building the code become convenient or necessary, we can look at splitting the necessary sub-packages
out into a public repo.

Thus this document exists largely either to sate your curiosity, or to serve as a guide for flooders to develop and release the cli :)

## Prerequisites

- go 1.9
- I recommend using `vg` for all go projects: https://github.com/GetStream/vg#installation.
- `cli` uses dep for dependency management:
  - install dep via brew `brew install dep`
  - or `go get -u github.com/golang/dep/cmd/dep`

## Getting started

```
vg init
vg ensure
make test
```
