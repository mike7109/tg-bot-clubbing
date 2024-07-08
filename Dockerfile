FROM golang:1.22-alpine AS build-env

RUN apk update && apk add --no-cache git openssh-client

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download -x

COPY . .
RUN sh .github/docker-build-app.sh

FROM alpine:latest

WORKDIR /
COPY --from=build-env /usr/src/app/build/app /app

CMD ["/app"]