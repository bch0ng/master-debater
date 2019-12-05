SUMMARY_NAME="$DOCKER_USERNAME/summary-service"
REMOVE_ALL_DOCKER_CONTAINERS="docker rm -vf \$(docker ps -a -q)"
docker pull $SUMMARY_NAME;

DEPLOY_SUMMARY="docker pull $SUMMARY_NAME; docker run -d --name micro-summary $SUMMARY_NAME"
#eval $REMOVE_ALL_DOCKER_CONTAINERS;
DOCKER_DELETE_PREV="docker rm -f micro-summary;"
eval $DOCKER_DELETE_PREV;
eval $DEPLOY_SUMMARY;
