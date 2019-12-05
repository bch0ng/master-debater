#!/bin/bash
# Variables
export DOCKER_BUILD_NAME="lleontan/gateway"
export ADDR=":443"
export TLSCERT="/etc/letsencrypt/live/api.sumsumsummary.me/fullchain.pem"
export TLSKEY="/etc/letsencrypt/live/api.sumsumsummary.me/privkey.pem"
export SESSIONKEY="sessionkey123"
export REDISADDR="redis://127.0.0.1:6379"
export DSN="your-password@tcp(127.0.0.1:3306)/psql_db"

docker login
MYSQL_ROOT_PASSWORD="your-password"
DB_NAME="lleontan/db"
GATEWAY_NAME="lleontan/gateway"
MESSAGING_NAME="lleontan/mongo-service"
SUMMARY_NAME="lleontan/summary-service"
REDIS_NAME="redis"
REMOVE_ALL_DOCKER_CONTAINERS="docker rm -vf \$(docker ps -a -q)"
ALTERNATIVE_DSN="root:your-password@tcp(psql_database:3306)/psql_db"
MESSAGESADDR="micro-messaging"
SUMMARYADDR="micro-summary"
DELETE_GATEWAY="docker rm -f info441-api"
DELETE_REDIS="docker rm -f my-redis"


LETS_ENCRYPT_MOUNT="-v /etc/letsencrypt:/etc/letsencrypt:ro"
LOGIN_MESSAGE="echo 'Login success!'"
CREATE_NETWORK="docker network create mynet"
DOCKER_LOGIN="echo 'Attempting docker login';docker login && echo 'Docker login success!'"
DEPLOY_DB="docker pull $DB_NAME; docker run -d --net mynet -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=psql_db --name psql_database $DB_NAME"
DEPLOY_REDIS="docker pull $REDIS_NAME; docker run -d --net mynet --name my-redis $REDIS_NAME"
MONGODB_INIT_SEED="mongo info441-a5 --eval 'db.createCollection(\"channel\") && db.createCollection(\"message\") && db.channel.createIndex( { \"email\" : 1 }, { unique : true } ) && db.channel.insertOne({ id: 1, name: \"general\", description: \"General channel\", private: false, members: [], createdAt: new Date(), creator: \"Admin\", editedAt: new Date()})'"
DEPLOY_MONGODB="docker run -d --net mynet --name my-mongodb mongo && docker exec my-mongodb bash -c \"$MONGODB_INIT_SEED\""
DEPLOY_MESSAGING="$DEPLOY_MONGODB; docker pull $MESSAGING_NAME; docker run -d --net mynet -e PORT=80 -e MONGO_URI=\"mongodb://my-mongodb/\" -e MONGO_DB=info441-a5 --expose 9000 --name micro-messaging $MESSAGING_NAME"
DEPLOY_SUMMARY="docker pull $SUMMARY_NAME; docker run -d --net mynet -e ADDR=:80 --expose 8000 --name micro-summary $SUMMARY_NAME"

DEPLOY_GATEWAY="docker pull $GATEWAY_NAME; docker run -d --net mynet --name info441-api -p 80:80 -p 443:443 -e ADDR=$ADDR -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e SESSIONKEY=$SESSIONKEY -e REDISADDR=$REDISADDR -e DSN=\"$ALTERNATIVE_DSN\" -e MESSAGESADDR=$MESSAGESADDR -e SUMMARYADDR=$SUMMARYADDR $GATEWAY_NAME"
eval "$DELETE_GATEWAY;$DELETE_REDIS; $DEPLOY_GATEWAY; $DEPLOY_REDIS;"
