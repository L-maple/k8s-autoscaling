#!/bin/bash


# hadoop benchmark脚本文件（手动）
#       
a=$1
b=$2
c=$3
curdir=`pwd`
# command: ./hdfs-TestDFSIO.sh  <write/read> <nrFiles> <size>
#       eg: ./hdfs-TestDFSIO.sh write 10 128
#       size: 单位MB
  
#----------------------------TestDFSIO----------------------------#

mkdir TestDFSIO_log

#for i in {1..$b}
for ((i=1; i<=$b; i++))
do
  hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -$a -nrFiles $b -size $c"MB" -resFile $curdir/TestDFSIO_log/TestDFSIO-$a-$i-$c.log
done