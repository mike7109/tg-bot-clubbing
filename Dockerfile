FROM golang:1.22-alpine AS build-env

RUN apk update && apk add --no-cache gcc musl-dev make

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download -x

COPY . .

ENV CGO_ENABLED=1

RUN sh .github/docker-build.sh

FROM alpine:latest

WORKDIR /
COPY --from=build-env /usr/src/app/build/app /app

CMD ["/app"]