#!/bin/sh
DB_BUILD_NAME="$DOCKER_USERNAME/debate-mongodb"
docker build -t $DB_BUILD_NAME .
docker login
docker push $DB_BUILD_NAME;
