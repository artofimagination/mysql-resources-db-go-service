#!/bin/bash
set -e

go build -i -race -ldflags "-X github.com/artofimagination/mysql-resources-db-go-service/config.AppVersion=$1" -o "$2"
chmod +x "$2"
