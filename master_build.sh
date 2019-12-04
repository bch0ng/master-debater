#!/bin/sh
(echo $DOCKER_USERNAME)
(cd servers/db/mysql && sh build.sh)
(cd servers/db/mongodb && sh build.sh)
(cd servers/gateway/ && sh build.sh)
(cd servers/debates/ && sh build.sh)
