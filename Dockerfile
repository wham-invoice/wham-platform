# Instructs the Docker builder what syntax to use when parsing the Dockerfile
# syntax=docker/dockerfile:1

# build stage
FROM golang:1.18-alpine as build

COPY . /platform
WORKDIR /platform

RUN go mod download

RUN go build -o /platform -buildvcs=false

# deploy stage
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /platform /platform

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/platform"]