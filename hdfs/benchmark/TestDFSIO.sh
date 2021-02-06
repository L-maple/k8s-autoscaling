#!/bin/bash


# hadoop benchmark脚本文件（手动）
#
nrFiles=$1
size=$2

curdir=`pwd`
# command: ./hdfs-TestDFSIO.sh  write-append <nrFiles> <size>
#       size: 单位MB
  
#----------------------------TestDFSIO----------------------------#

mkdir testdfsio_log
hadoop jar "${HADOOP_HOME}"/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -write -nrFiles 1 -size $size"MB" -resFile "$curdir"/testdfsio_log/TestDFSIO-write-$nrFiles-1-$size.log
for ((i=1; i<"$nrFiles"; i++))
do
  hadoop jar "${HADOOP_HOME}"/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -append -nrFiles 1 -size $size"MB" -resFile "$curdir"/testdfsio_log/TestDFSIO-append-$nrFiles-$i-$size.log
done
# hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -clean

