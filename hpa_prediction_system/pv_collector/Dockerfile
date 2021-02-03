FROM centos

WORKDIR /usr/bin

COPY . .

# the pv_collector depends on the `df` and `iostat` command
RUN echo "67.219.148.138 mirrorlist.centos.org" >> /etc/hosts
RUN yum install -y sysstat

ENTRYPOINT ["./pv_collector"]
