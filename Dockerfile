FROM golang:1-alpine as dev
WORKDIR /app
COPY src .
COPY ./conf/runner.conf /runner.conf
RUN apk add --update git gcc libc-dev npm && \
    go get github.com/pilu/fresh && \
    chown -R 1000. /app

FROM golang:1-alpine as gobuild
WORKDIR /app
COPY src .
RUN apk add --update git gcc libc-dev && \
    CGO_ENABLED=1 GOOS=linux go build -o shorty -ldflags="-s -w" .

FROM node:lts-alpine as assetbuild
WORKDIR /app
COPY src/assets .
RUN npm install && npm run build

FROM alpine:latest as prod
WORKDIR /app
VOLUME /data
COPY --from=gobuild /app/shorty .
COPY --from=gobuild /app/templates ./templates
COPY --from=assetbuild /app/dist ./assets/dist
CMD ["/app/shorty"]
