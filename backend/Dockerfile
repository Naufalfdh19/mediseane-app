# syntax=docker/dockerfile:1

FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./ . ./
RUN GOPROXY="https://goproxy.io,direct" go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend

FROM alpine AS build-release-stage

WORKDIR /

COPY --from=build-stage /backend /backend
COPY .env ./

EXPOSE 8000

ENTRYPOINT ["/backend"]
