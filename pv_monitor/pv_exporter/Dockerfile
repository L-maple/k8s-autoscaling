FROM golang:1.14.3-alpine AS build

WORKDIR /src
COPY . .

RUN cp -R .kube /root

ENTRYPOINT ["./pv_exporter", "-statefulset", "ubuntu", "-namespace", "default"]
