DOCKER_BUILD_NAME="$DOCKER_USERNAME/summary-service"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o summary .
echo "DOCKER BUILD NAME: $DOCKER_BUILD_NAME"
docker build -t $DOCKER_BUILD_NAME .
docker login
docker push $DOCKER_BUILD_NAME;
go clean
