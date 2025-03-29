FROM node:20-alpine AS web

WORKDIR /web
COPY ./web/package.json package.json 
COPY ./web/package-lock.json package-lock.json
RUN npm install
COPY ./web /web
RUN npm run build

FROM golang:1.22 AS builder
ENV             GO11MODULE=on
WORKDIR         /go/
COPY            go.mod go.sum ./
RUN             go mod download
COPY            . ./
COPY --from=web /web/dist ./web/dist
RUN             make _docker_install

FROM scratch
ARG NAME=survey-repository
WORKDIR /app
COPY --from=builder /go/${NAME} .
