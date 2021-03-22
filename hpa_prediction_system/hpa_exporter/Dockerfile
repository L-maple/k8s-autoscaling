FROM centos

WORKDIR /usr/bin

COPY . .

EXPOSE 30001 30002 30003

ENTRYPOINT ["./hpa_exporter", "-statefulset", "hdfs-datanode", "-namespace", "monitoring"]
#ENTRYPOINT ["./hpa_exporter", "-statefulset", "advanced-tidb-tikv", "-namespace", "tidb-cluster"]
