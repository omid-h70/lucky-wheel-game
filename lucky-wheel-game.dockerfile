# Base go image
# Build in Vendor mod
#
#
#FROM golang:1.18 as builder
FROM golang:1.18-alpine

ARG BINARY_PATH=/app/cmd/api

RUN mkdir /app

COPY . /app

WORKDIR $BINARY_PATH

RUN CGO_ENABLED=0 go build -mod=vendor -o gameService

RUN chmod +x gameService

CMD [ "./gameService" ]

###### Production Image - tiny One !
#
#FROM alpine:latest

#RUN mkdir /gameService

#COPY . /gameService

#WORKDIR $BINARY_PATH

#RUN CGO_ENABLED=0 go build -mod=vendor -o gameService

#RUN chmod +x $BINARY_PATH/gameService

#CMD [ "$BINARY_PATH" ]