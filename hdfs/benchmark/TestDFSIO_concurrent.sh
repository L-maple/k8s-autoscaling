#!/bin/bash

# hadoop benchmark脚本文件（手动）
nrFiles=$1
concurrent=$2
size=$3

curdir=`pwd`
# command: ./hdfs-TestDFSIO.sh  write-append <nrFiles> <concurrent> <size>
#       size: 单位MB

#----------------------------TestDFSIO----------------------------#

mkdir testdfsio_log_"$nrFiles"_"$size"
hadoop jar "${HADOOP_HOME}"/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -write -nrFiles "$concurrent" -size $size"MB" -resFile "$curdir"/testdfsio_log_"$nrFiles"_"$size"/TestDFSIO-write-$nrFiles-1-$size.log
for ((i=1; i<"$nrFiles"; i++))
do
  hadoop jar "${HADOOP_HOME}"/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -append -nrFiles "$concurrent" -size $size"MB" -resFile "$curdir"/testdfsio_log_"$nrFiles"_"$size"/TestDFSIO-append-$nrFiles-$i-$size.log
done