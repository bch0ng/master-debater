#!/bin/bash
# Variables
export DOCKER_BUILD_NAME="bch0ng/gateway"
export ADDR=":443"
export TLSCERT="/etc/letsencrypt/live/api.sumsumsummary.me/fullchain.pem"
export TLSKEY="/etc/letsencrypt/live/api.sumsumsummary.me/privkey.pem"
export SESSIONKEY="sessionkey123"
export REDISADDR="redis://127.0.0.1:6379"
export DSN="your-password@tcp(127.0.0.1:3306)/psql_db"

bash build.sh

docker login
docker push $DOCKER_BUILD_NAME;

DOCKER_COMMANDS="echo 'Login success!'; docker rm -f info441-api my-redis && docker pull $DOCKER_BUILD_NAME && docker run -d --name info441-api -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e ADDR=$ADDR -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e SESSIONKEY=$SESSIONKEY -e REDISADDR=$REDISADDR -e DSN=\"$DSN\" $DOCKER_BUILD_NAME && docker run -d --name my-redis redis;"

echo 'Attempting default INFO441 SSH config info441-client...'
#reminder, Leon's username is bch0ng
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
