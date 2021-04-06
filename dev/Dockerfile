FROM golang:1.15

ARG SERVICE_NAME=mysql-resources-db-go-service

RUN go get github.com/githubnemo/CompileDaemon
RUN go get github.com/go-delve/delve/cmd/dlv

RUN mkdir -p /go/src/github.com/artofimagination/$SERVICE_NAME
WORKDIR /go/src/github.com/artofimagination/$SERVICE_NAME

VOLUME /go/src/github.com/artofimagination/$SERVICE_NAME

COPY ./entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["run"]

EXPOSE 80
