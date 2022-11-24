ARG alpine_version=3.16
ARG go_version=1.18.2

# Server binary builder
FROM golang:${go_version}-alpine${alpine_version} as base

ARG SSH_PRIVATE_KEY
ARG KEN_CRT
ARG CI_PROJECT_NAME

RUN apk add --no-cache git
RUN apk add --update --no-cache curl build-base

RUN git config --global "url.ssh://git@gitlab.kenda.com.tw:4222".insteadOf "https://gitlab.kenda.com.tw"

RUN apk update && apk add openssh

ARG SSH_PUBLIC_KEY

RUN mkdir ~/.ssh && \
    echo "${SSH_PRIVATE_KEY}" > ~/.ssh/id_rsa && \
    echo "${SSH_PUBLIC_KEY}" > ~/.ssh/id_rsa.pub && \
    chmod 600 ~/.ssh/id_rsa && \
    chmod 600 ~/.ssh/id_rsa.pub && \
    ssh-keyscan -Ht ecdsa -p 4222 gitlab.kenda.com.tw,192.1.1.159 >> ~/.ssh/known_hosts

RUN echo "${KEN_CRT}" >> /etc/ssl/certs/ca-certificates.crt

ENV REPO_DIR ${GOPATH}/src/gitlab.com.kenda.com.tw/kenda/${CI_PROJECT_NAME}

FROM base AS test-builder
ARG CI_PROJECT_NAME

COPY ./ ${REPO_DIR}
WORKDIR ${REPO_DIR}

ENV GOPRIVATE *.kenda.com.tw
RUN go mod download
RUN go vet ./...
RUN go test -race -coverprofile .testCoverage.txt ./...
RUN go tool cover -func .testCoverage.txt
RUN go build -race -ldflags "-extldflags '-static'" -o /opt/${CI_PROJECT_NAME}/server

# Deployable with server binary and UI dist
FROM alpine:${alpine_version}
ARG CI_PROJECT_NAME

WORKDIR /root/
COPY --from=test-builder /opt/${CI_PROJECT_NAME}/server /root/server
CMD ["/bin/sh"]
