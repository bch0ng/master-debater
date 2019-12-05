# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server
FROM alpine:latest
RUN apk add --no-cache bash ca-certificates
COPY gateway /gateway
EXPOSE 443:443
ENTRYPOINT ["/gateway"]