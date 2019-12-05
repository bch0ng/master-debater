#!/bin/bash
DB_BUILD_NAME="$DOCKER_USERNAME/db"
echo "DOCKER BUILD NAME: $DB_BUILD_NAME"
docker build -t $DB_BUILD_NAME .
docker login
docker push $DB_BUILD_NAME;
