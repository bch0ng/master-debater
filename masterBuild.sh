#runs all the build scripts in the repo
(cd clients/summary/; sh build.sh)
(cd servers/db/; sh buildDbDocker.sh)
(cd servers/gateway/; sh build.sh)
(cd servers/messaging/; sh buildMongoService.sh)
(cd servers/summary/; sh buildSummaryMicroservice.sh)
