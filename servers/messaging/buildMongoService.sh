docker build -t $DOCKER_USERNAME/mongo-service .
docker login
docker push $DOCKER_USERNAME/mongo-service;
