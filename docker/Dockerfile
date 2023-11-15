FROM golang:1.21.3-alpine as gobuff

EXPOSE 8080

ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache --upgrade apk-tools \
    && apk add --no-cache bash curl openssl git build-base nodejs npm sqlite sqlite-dev mysql-client vim postgresql libpq postgresql-contrib libc6-compat

# Installing linter
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.55.2

# Installing Yarn
RUN npm i -g --no-progress yarn \
    && yarn config set yarn-offline-mirror /npm-packages-offline-cache \
    && yarn config set yarn-offline-mirror-pruning true

# Installing buffalo binary
RUN go install github.com/gobuffalo/cli/cmd/buffalo@latest
RUN go install github.com/gobuffalo/buffalo-pop/v3@latest

RUN mkdir /src
WORKDIR /src

FROM gobuff as builder

ARG version=0.0.0
ARG prerelease
ARG githash=""
ARG buildstamp=""

RUN mkdir /app
WORKDIR /app
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go version

ADD . .
RUN buffalo build --static -o /bin/app --ldflags "-X github.com/YaleSpinup/proctor/proctor.Version=$version -X github.com/YaleSpinup/proctor/proctor.VersionPrerelease=$prerelease -X github.com/YaleSpinup/proctor/proctor.GitHash=$githash -X github.com/YaleSpinup/proctor/proctor.BuildStamp=$buildstamp"

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache python3
#RUN ln -s /usr/bin/python3 /usr/bin/python
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache wget gettext && \
    wget -nv "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" && \
    unzip awscli-bundle.zip && \
    ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws && \
    /usr/local/bin/aws --version

WORKDIR /bin/

COPY --from=builder /bin/app .
COPY docker/import_config.sh /bin

ENV GO_ENV=production
ENV PORT=8080
ENV ADDR=0.0.0.0

EXPOSE 8080

CMD /bin/import_config.sh && /bin/app
