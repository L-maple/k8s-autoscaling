FROM golang:1.14.3-alpine AS build

WORKDIR /src
COPY . .

RUN go build -o pv_exporter .
ENTRYPOINT ["./pv_exporter"]
