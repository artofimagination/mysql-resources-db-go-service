FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/mysql-resources-db-go-service

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++
RUN go mod tidy
RUN cd $GOPATH/src/mysql-resources-db-go-service/ && go build main.go

# This container exposes port 8080 to the outside world
EXPOSE 8080

RUN chmod 0766 $GOPATH/src/mysql-resources-db-go-service/scripts/init.sh

# Run the executable
CMD ["./scripts/init.sh"]