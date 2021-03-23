FROM centos:latest

WORKDIR /usr/bin

COPY . .

# the pv_collector depends on the `df` and `iostat` command
RUN echo "nameserver 10.21.1.205" >> /etc/resolv.conf
RUN echo "nameserver 8.8.8.8" >> /etc/resolv.conf
RUN echo "nameserver 8.8.4.4" >> /etc/resolv.conf
RUN yum install -y sysstat && yum install -y lvm2

ENTRYPOINT ["./pv_collector", "-serverAddress", "kubernetes://hdfs-hpa-exporter-service.monitoring:30002/"]
#ENTRYPOINT ["./pv_collector", "-serverAddress", "kubernetes://tidb-hpa-exporter-service.monitoring:30002/"]
