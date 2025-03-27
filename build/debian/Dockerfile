FROM golang:1.22 AS builder
ENV             GO11MODULE=on
WORKDIR         /go/
COPY            go.mod go.sum ./
RUN             go mod download
COPY            . ./
#COPY --from=web /web/dist ./web/dist
RUN             make _docker_install

FROM scratch
ARG NAME=survey-repository
WORKDIR /app
COPY --from=builder /go/${NAME} .
