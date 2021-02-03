FROM centos

WORKDIR /usr/bin

COPY . .

# the pv_collector depends on the `df` and `iostat` command
RUN rm -f /etc/resolv.conf

COPY /etc/resolv.conf /etc/

RUN yum install -y sysstat

ENTRYPOINT ["./pv_collector"]
