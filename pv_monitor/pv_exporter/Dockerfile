FROM golang:1.14.3-alpine AS build

WORKDIR /src
COPY . .

RUN cp pv_exporter /usr/local/bin &&\
    cp pv_exporter /usr/bin

ENTRYPOINT ["./pv_exporter", "-statefulset", "ubuntu"]
