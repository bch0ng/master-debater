#Deploys the gateway on the server
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
export ADDR=":443"
export TLSCERT="/etc/letsencrypt/live/api.sumsumsummary.me/fullchain.pem"
export TLSKEY="/etc/letsencrypt/live/api.sumsumsummary.me/privkey.pem"
export SESSIONKEY="sessionkey123"
export REDISADDR="redis://127.0.0.1:6379"
export DSN="your-password@tcp(127.0.0.1:3306)/psql_db"

REMOVE_PREV_API="docker rm -f info441-api"
DEPLOY_GATEWAY="$REMOVE_PREV_API; docker pull $GATEWAY_NAME; docker run -d --net mynet --name info441-api -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e ADDR=$ADDR -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e SESSIONKEY=$SESSIONKEY -e REDISADDR=$REDISADDR -e DSN=\"$ALTERNATIVE_DSN\" -e MESSAGESADDR=$MESSAGESADDR -e SUMMARYADDR=$SUMMARYADDR $GATEWAY_NAME"

DEFAULT_PRIV_KEY="$HOME/.ssh/id_rsa"

ssh -i $DEFAULT_PRIV_KEY bch0ng@ec2-54-218-179-6.us-west-2.compute.amazonaws.com $DEPLOY_GATEWAY exit; [ $? -eq 255 ]
