# Build node app
FROM node:12.14.1-alpine as build-node

COPY ./web-src /web-src

RUN cd /web-src && \
    apk add --no-cache python make g++ && \
	npm ci && \
	npm run build

# Golang build
FROM golang:1.13.7-buster as build-go

COPY . /app

RUN cd /app && \
    go get -d -v ./... && \
    go build cmd/server.go

# All together now :)
FROM alpine:3.9.5

RUN mkdir -p /opt/brunel/web/ && \
    apk add --no-cache libc6-compat

WORKDIR /opt/brunel/

COPY --from=build-go /app/server /opt/brunel/server
COPY --from=build-node /web-src/build/ /opt/brunel/web/
COPY ./docker-server.yaml /opt/brunel/brunel.yaml

ENTRYPOINT [ "/opt/brunel/server" ]
