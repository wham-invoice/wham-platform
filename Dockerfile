# Instructs the Docker builder what syntax to use when parsing the Dockerfile
# syntax=docker/dockerfile:1

# build stage
FROM golang:1.18-alpine as build

COPY . /wham-platform
WORKDIR /wham-platform

RUN go mod download

RUN go build -o /wham-platform-bin -buildvcs=false

# deploy stage
FROM golang:1.18-alpine

WORKDIR /

COPY --from=build /wham-platform-bin /wham-platform-bin

EXPOSE 8080

ENTRYPOINT ["/wham-platform-bin"]