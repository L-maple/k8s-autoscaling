FROM centos

WORKDIR /usr/bin

COPY . .

# the pv_collector depends on the `df` and `iostat` command
RUN firewall-cmd --permanent --zone=trusted --add-interface=docker0
RUN firewall-cmd --reload

RUN yum install -y sysstat

ENTRYPOINT ["./pv_collector"]
