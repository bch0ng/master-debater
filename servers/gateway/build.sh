#!/bin/bash
DOCKER_BUILD_NAME="$DOCKER_USERNAME/gateway"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o gateway .
echo "DOCKER BUILD NAME: $DOCKER_BUILD_NAME"
docker build -t $DOCKER_BUILD_NAME .
docker login
docker push $DOCKER_BUILD_NAME;
go clean
