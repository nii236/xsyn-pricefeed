# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.3-buster AS build
ADD https://www.google.com /time.now

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ethusd ./ethusd
COPY ./middleware ./middleware
COPY ./supseth ./supseth
COPY *.go ./

RUN mkdir ./dist
RUN GOOS=linux GOARCH=amd64 go build -o ./dist/xsyn-pricefeed

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=build /app/dist/xsyn-pricefeed /xsyn-pricefeed
EXPOSE 8080

USER nonroot:nonroot

CMD ["/xsyn-pricefeed"]
