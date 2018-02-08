FROM golang:1.9-alpine

ARG GITHUB_TOKEN
RUN \
      apk add --no-cache curl git bash make &&\
      echo installing dep &&\
      mkdir -p /tmp/dockerfile &&\
      cd /tmp/dockerfile &&\
      \
      curl -fsSL "https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64" > dep &&\
      GODEP_SHA256=322152b8b50b26e5e3a7f6ebaeb75d9c11a747e64bbfd0d8bb1f4d89a031c2b5 \
        echo "$GODEP_SHA256 dep" | sha256sum &&\
      chmod 0755 dep &&\
      mv dep /usr/bin/dep &&\
      \
      echo installing goreleaser &&\
      curl -fsSL "https://github.com/goreleaser/goreleaser/releases/download/v0.46.4/goreleaser_Linux_x86_64.tar.gz" > goreleaser.tar.gz &&\
      GORELEASER_SHA256=2f3144de881e8204dcbd5c561eb6779bc8b451bddcf977c4c58819f0cce1b670 \
        echo $GORELEASER_SHA256 goreleaser.tar.gz | sha256sum &&\
      tar -xf goreleaser.tar.gz &&\
      mv goreleaser /usr/bin/goreleaser &&\
      \
      echo cleaning up &&\
      apk del curl &&\
      cd / &&\
      rm -rf /tmp/dockerfile &&\
    git config --global url."https://github.com".insteadOf git://github.com &&\
    git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# for tests:
RUN go get github.com/stretchr/testify/assert

WORKDIR /go/src/github.com/flood-io/cli
ADD Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

ADD . ./
ARG GIT_SHA
