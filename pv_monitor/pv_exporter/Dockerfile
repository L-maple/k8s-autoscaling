FROM golang:1.14.3-alpine AS build

WORKDIR /usr/bin

COPY . .

ENTRYPOINT ["./pv_exporter", "-statefulset", "ubuntu", "-namespace", "default", "-incluster", "true"]
