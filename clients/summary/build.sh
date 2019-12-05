echo "Attempting client build."
DOCKER_BUILD_NAME="$DOCKER_USERNAME/summary"
echo "DOCKER Client BUILD NAME: $DOCKER_BUILD_NAME"
docker build -t $DOCKER_BUILD_NAME .
