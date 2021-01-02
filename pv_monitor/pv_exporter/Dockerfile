FROM golang:1.14.3-alpine AS build

WORKDIR /usr/bin

COPY . .

EXPOSE 30001 30002 30003

ENTRYPOINT ["./pv_exporter", "-statefulset", "ubuntu", "-namespace", "default"]
