FROM golang:1.17-buster as dev

ARG GolangCILintVersion=1.42.1

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
        sh -s -- -b $(go env GOPATH)/bin v${GolangCILintVersion} \
    && go get github.com/smartystreets/goconvey