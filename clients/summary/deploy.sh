#!/bin/sh
# Variables
export DOCKER_BUILD_NAME="$DOCKER_USERNAME/summary"
OTHER_BUILD_NAME="clientserver"
sh build.sh

docker login
docker push $DOCKER_BUILD_NAME

DOCKER_STOP_ALL="docker container stop \$(docker container ls -aq);
                  docker container rm \$(docker container ls -aq)"
DOCKER_STOP_SPECIFIC="sudo docker rm -f $DOCKER_BUILD_NAME"
DOCKER_COPY_SERVER_STATIC_TO_DOCKER="-v /usr/share/nginx/html:/usr/share/nginx/html:ro"
DOCKER_COPY_CERTS_FOLDER="-v /etc/letsencrypt:/etc/letsencrypt:ro"
DOCKER_MOUNT_CONF="-v /etc/nginx/conf.d/:/etc/nginx/conf.d/:ro"
DOCKER_COMMANDS="echo 'Login success!';
    $DOCKER_STOP_ALL;
    docker pull $DOCKER_BUILD_NAME;
    sudo docker run -d $DOCKER_COPY_CERTS_FOLDER -p 80:80 -p 443:443 --name $OTHER_BUILD_NAME $DOCKER_USERNAME/summary"

#Default private key is the key for your laptop's private key
DEFAULT_PRIV_KEY="~/.ssh/id_rsa"
echo 'Attempting default INFO441 SSH config info441-client...'
if ssh -t info441-client $DOCKER_COMMANDS exit; [ $? -eq 255 ]
then
    echo 'Unsuccessful... Trying manual login.'
    #Ask for ssh login username
    echo 'sumsumsummary.me SSH login username:'
    read USERNAME
    # Ask user for location of their privKey
    #/etc/letsencrypt/live/sumsumsummary.me/cert.pem for cert
    #/etc/letsencrypt/live/sumsumsummary.me/privkey.pem for key
    #echo 'Location of sumsumsummary.me private key:'
    #read PRIV_KEY_LOC

    #if
    ssh -i $DEFAULT_PRIV_KEY $USERNAME@sumsumsummary.me $DOCKER_COMMANDS; [ $? -eq 255 ]
    #then
    #    echo 'Unsuccessful... Exiting.'
    #fi
fi
