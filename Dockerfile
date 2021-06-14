FROM golang:latest AS build-stage

COPY . /opt/wallet-service
WORKDIR /opt/wallet-service

ENV GO111MODULE on

# build
RUN go build . -o wallet-service

# deploy
FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /opt/wallet-service
COPY --from=build-stage /opt/wallet-service/wallet-service  ./wallet-service

ENTRYPOINT [ "/opt/wallet-service/wallet-service" ]
CMD [ "run" ]
