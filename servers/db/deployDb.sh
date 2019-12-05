DOCKER_BUILD_NAME="psql_database"
MYSQL_ROOT_PASSWORD="your-password"
DOCKER_STOP_SPECIFIC="sudo docker rm -f $DOCKER_BUILD_NAME"
DOCKER_STOP_ALL="docker container stop \$(docker container ls -aq)
                  docker container rm \$(docker container ls -aq);"
docker rm -f $DOCKER_BUILD_NAME;
docker run -d \
-p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=psql_db \
--name $DOCKER_BUILD_NAME $DOCKER_USERNAME/db;
