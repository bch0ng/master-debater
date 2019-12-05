#Build Everything
bash masterBuild.sh


#run the summary client
(cd clients/summary/; sh deploy.sh)

#Deploy the gateway server with the microservices
(cd servers/; sh deployServer.sh)
