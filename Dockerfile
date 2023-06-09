# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.3-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ethusd ./ethusd
COPY ./middleware ./middleware
COPY ./supseth ./supseth
COPY ./erc20 ./erc20
COPY *.go ./

RUN mkdir ./dist
RUN GOOS=linux GOARCH=amd64 go build -o ./dist/xsyn-pricefeed

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=build /app/dist/xsyn-pricefeed /xsyn-pricefeed
COPY ./CHECKS /
USER nonroot:nonroot

ENTRYPOINT ["/xsyn-pricefeed", "serve"]
