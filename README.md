# `flood` CLI

The Flood Command Line Interface

## Usage

```bash
flood help
```

### Authenticating

Before you can do anything, we need to authenticate your local machine with your
Flood account:

    $ flood login

    What's your username:
    Enter a value: user@exampe.com

    What's your password (masked):
    Enter a value: **************************************************

    Welcome back Ivan Vanderbyl!

This will store a temporary authentication token on your machine.


You can deregister this machine by running `flood logout` at any time.

## Development

- protoc - https://github.com/google/protobuf/releases
  - e.g. `brew install protoc`
  - or e.g. https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
- gogo/protobuf
  - go get github.com/gogo/protobuf/proto
  - go get github.com/gogo/protobuf/protoc-gen-gofast
  - go get github.com/gogo/protobuf/gogoproto

