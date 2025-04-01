FROM node:20-alpine AS web

WORKDIR /web
COPY ./web/package.json package.json 
COPY ./web/package-lock.json package-lock.json
RUN npm install
COPY ./web /web
RUN npm run build

FROM golang:1.22 AS builder
ENV  GO11MODULE=on
WORKDIR /go/
# Install nfpm to build rpm
RUN  echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list && \
    apt update && \
    apt install nfpm
COPY go.mod go.sum ./
RUN  go mod download
COPY . ./
COPY --from=web /web/dist ./web/dist
RUN  make _docker_install
RUN  ./build/scripts/rpm.sh

# Final with only built files
FROM scratch
ARG NAME=survey-repository
COPY --from=builder /go/${NAME} .
COPY --from=builder /go/*.rpm .
