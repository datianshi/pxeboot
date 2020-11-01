FROM golang:1.14.3 AS build-env

ADD . /app
WORKDIR /app

RUN go build -o /pxeboot cmd/pxe/main.go 

FROM ubuntu:latest
COPY --from=build-env /pxeboot /

CMD [ "/pxeboot", "-config", "/data/config.yaml" ]