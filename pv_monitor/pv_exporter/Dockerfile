FROM centos

WORKDIR /usr/bin

COPY . .

EXPOSE 30001 30002 30003

ENTRYPOINT ["./pv_exporter", "-statefulset", "ubuntu", "-namespace", "default"]
