#!/bin/sh
docker build -t $DOCKER_USERNAME/micro-debates .
docker login
docker push $DOCKER_USERNAME/micro-debates;
