#!/bin/bash
# Variables

export ADDR=":443"
export TLSCERT="/etc/letsencrypt/live/api.sumsumsummary.me/fullchain.pem"
export TLSKEY="/etc/letsencrypt/live/api.sumsumsummary.me/privkey.pem"
export SESSIONKEY="sessionkey123"
export REDISADDR="my-redis:6379"
export RABBITMQADDR="amqp://my-rabbitmq"
export DSN="your-password@tcp(127.0.0.1:3306)/psql_db"

MYSQL_ROOT_PASSWORD="your-password"
DB_NAME="$DOCKER_USERNAME/db"
GATEWAY_NAME="$DOCKER_USERNAME/gateway"
MESSAGING_NAME="$DOCKER_USERNAME/mongo-service"
SUMMARY_NAME="$DOCKER_USERNAME/summary-service"
REDIS_NAME="redis"
REMOVE_ALL_DOCKER_CONTAINERS="docker rm -vf \$(docker ps -a -q)"
ALTERNATIVE_DSN="root:your-password@tcp(psql_database:3306)/psql_db"
MESSAGESADDR="micro-messaging"
SUMMARYADDR="micro-summary"
#TODO, REDIS needs a dockerfile and build script

DEPLOY_RABBITMQ="docker pull rabbitmq;docker run -d --net mynet --name my-rabbitmq rabbitmq"
LOGIN_MESSAGE="echo 'Login success!'"
CREATE_NETWORK="docker network create mynet"
DOCKER_LOGIN="echo 'Attempting docker login';docker login && echo 'Docker login success!'"
DEPLOY_DB="docker pull $DB_NAME; docker run -d --net mynet -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=psql_db --name psql_database $DB_NAME"
DEPLOY_REDIS="docker pull $REDIS_NAME; docker run -d --net mynet --name my-redis $REDIS_NAME"
MONGODB_INIT_SEED="mongo info441-a5 --eval 'db.createCollection(\"channel\") && db.createCollection(\"message\") && db.channel.createIndex( { \"email\" : 1 }, { unique : true } ) && db.channel.insertOne({ id: 1, name: \"general\", description: \"General channel\", private: false, members: [], createdAt: new Date(), creator: \"Admin\", editedAt: new Date()})'"
DEPLOY_MONGODB="docker run -d --net mynet --name my-mongodb mongo && docker exec my-mongodb bash -c \"$MONGODB_INIT_SEED\""
DEPLOY_GATEWAY="docker pull $GATEWAY_NAME; docker run -d --net mynet --name info441-api -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e ADDR=$ADDR -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e SESSIONKEY=$SESSIONKEY -e REDISADDR=$REDISADDR -e DSN=\"$ALTERNATIVE_DSN\" -e MESSAGESADDR=$MESSAGESADDR -e SUMMARYADDR=$SUMMARYADDR -e RABBITMQADDR=$RABBITMQADDR $GATEWAY_NAME"
DEPLOY_MESSAGING="$DEPLOY_MONGODB; docker pull $MESSAGING_NAME; docker run -d --net mynet -e PORT=80 -e MONGO_URI=\"mongodb://my-mongodb/\" -e MONGO_DB=info441-a5 -e RABBITMQADDR=amqp://my-rabbitmq --expose 9000 --name micro-messaging $MESSAGING_NAME"
DEPLOY_SUMMARY="docker pull $SUMMARY_NAME; docker run -d --net mynet -e ADDR=:80 --expose 8000 --name micro-summary $SUMMARY_NAME"

DOCKER_COMMANDS="$LOGIN_MESSAGE;$REMOVE_ALL_DOCKER_CONTAINERS;$CREATE_NETWORK;$DEPLOY_RABBITMQ;echo 'waiting!'; sleep 30s;$DEPLOY_MESSAGING; $DEPLOY_SUMMARY;$DEPLOY_DB; echo 'waiting!'; sleep 30s; $DEPLOY_GATEWAY; $DEPLOY_REDIS;"
DOCKER_DEPLOY_DB="docker rm -f $DOCKER_USERNAME/db && docker pull $DOCKER_USERNAME/db && docker run -d --name psql_database -p 3306:3306 $DOCKER_USERNAME/db;"


echo 'Attempting default INFO441 SSH config info441-API...\n'
#reminder, Leon's username is lleontan
DEFAULT_PRIV_KEY="$HOME/.ssh/id_rsa"
if ssh -t info441-api $DOCKER_COMMANDS exit; [ $? -eq 255 ]
then
  echo 'Unsuccessful... Trying manual login.'
  #Ask for ssh login username
  echo 'sumsumsummary.me SSH login username:'
  read USERNAME
  # Ask user for location of their privKey
  echo 'Location of sumsumsummary.me private key or default_location:'
  read PRIV_KEY_LOC
  if [ $PRIV_KEY_LOC == "default_location" ]
  then PRIV_KEY_LOC=$DEFAULT_PRIV_KEY
  fi
  echo $PRIV_KEY_LOC
  if ssh -i $PRIV_KEY_LOC $USERNAME@ec2-54-218-179-6.us-west-2.compute.amazonaws.com $DOCKER_COMMANDS exit; [ $? -eq 255 ]
  then
    echo 'Unsuccessful... Exiting.'
  fi
fi
echo "Server deploy complete"
