ARG EXECUTABLE_NAME=mysql-resources-db-go-service

FROM golang:1.15.2-alpine AS builder

ENV ROOT_PACKAGE=github.com/artofimagination/$EXECUTABLE_NAME

ADD . $GOPATH/src/$ROOT_PACKAGE
WORKDIR $GOPATH/src/$ROOT_PACKAGE

RUN apk add --update g++
RUN go mod tidy
RUN go build -ldflags "-X $ROOT_PACKAGE/config.AppVersion=$APP_VERSION" main.go

FROM alpine:latest

ARG EXECUTABLE_NAME

RUN set -eux; \
  apk add --no-cache ca-certificates

WORKDIR /usr/local/bin/

COPY --from=builder /tmp/$EXECUTABLE_NAME ./$EXECUTABLE_NAME
RUN chmod 0766 $GOPATH/src/$ROOT_PACKAGE/scripts/init.sh

EXPOSE 8080

# Run the executable
CMD ["./scripts/init.sh"]