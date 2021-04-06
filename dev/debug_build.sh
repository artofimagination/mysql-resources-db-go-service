#!/bin/bash
set -e

go build -i -race -gcflags "all=-N -l" -ldflags "-X github.com/artofimagination/mysql-resources-db-go-service/config.AppVersion=$1" -o "$2"
chmod +x "$2"
