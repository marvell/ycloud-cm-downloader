# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.217.1/containers/go/.devcontainer/base.Dockerfile

ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN set -xe; \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin