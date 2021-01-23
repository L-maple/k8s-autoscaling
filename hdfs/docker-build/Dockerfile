FROM       ubuntu:16.04
MAINTAINER nwpu

# finish the set about ssh
RUN apt-get update \
 && apt-get install -y vim \
 && apt install -y net-tools \
 && apt install -y iputils-ping \
 && apt install -y vim \
 && echo 'root:root' |chpasswd \
 && apt-get install -y openssh-server \
 && sed -ri 's/^PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config \
 && sed -ri 's/#   StrictHostKeyChecking ask/StrictHostKeyChecking no/' /etc/ssh/ssh_config \
 && mkdir /root/.ssh \
 && ssh-keygen -t rsa -P "" -f /root/.ssh/id_rsa \
 && cat /root/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys
RUN sed -ri 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config
RUN sed -ri 's/Port 22/Port 31822/g' /etc/ssh/sshd_config

# copy the src of java and hadoop
RUN mkdir /custom
RUN mkdir /custom/hdfs_dir
RUN mkdir /custom/prometheus_export
ADD hadoop-2.9.0.tar.gz /custom/
ADD jdk-8u11-linux-x64.tar.gz /custom/
COPY prometheus_export /custom/prometheus_export/
COPY conf /custom/hadoop-2.9.0/etc/hadoop/
COPY jar /custom/hadoop-2.9.0/share/hadoop/hdfs/
COPY jar_lib /custom/hadoop-2.9.0/share/hadoop/hdfs/lib/
COPY start.sh /custom
RUN sed -i '25c export JAVA_HOME=/custom/jdk1.8.0_11' /custom/hadoop-2.9.0/etc/hadoop/hadoop-env.sh
RUN echo "export HADOOP_SSH_OPTS=\"-p 31822\"" >> /custom/hadoop-2.9.0/etc/hadoop/hadoop-env.sh

# set the envionment variables of java and hadoop
ENV JAVA_HOME=/custom/jdk1.8.0_11 \
    JRE_HOME=/custom/jdk1.8.0_11/jre \
    HADOOP_HOME=/custom/hadoop-2.9.0 \
    CLASSPATH=./:/custom/jdk1.8.0_11/lib:/custom/jdk1.8.0_11/jre/lib \
    PATH=$PATH:/custom/hadoop-2.9.0/bin:/custom/hadoop-2.9.0/sbin:/custom/jdk1.8.0_11/bin

# Hdfs ports
EXPOSE 9000 50010 50020 50070 50075 50090
# Mapred ports
EXPOSE 19888
# Yarn ports
EXPOSE 8030 8031 8032 8033 8040 8042 8088
# Other ports
EXPOSE 49707 31822
# Prometheus port
EXPOSE 10318

RUN chmod 777 /custom/start.sh
CMD ["/custom/start.sh"]
