docker stop micro-messaging && docker rm micro-messaging
docker run -d --net mynet --name my-rabbitmq rabbitmq
