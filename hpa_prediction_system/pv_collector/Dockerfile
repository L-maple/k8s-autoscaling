FROM centos

WORKDIR /usr/bin

COPY . .

# the pv_collector depends on the `df` and `iostat` command
RUN yum install -y sysstat

ENTRYPOINT ["./pv_collector"]
