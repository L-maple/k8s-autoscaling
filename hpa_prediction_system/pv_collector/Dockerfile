FROM centos

WORKDIR /usr/bin

COPY . .

# the pv_collector depends on the `df` and `iostat` command
RUN echo "nameserver 10.21.1.205" >> /etc/resolv.conf
RUN yum install -y sysstat

ENTRYPOINT ["./pv_collector"]
