#!/bin/bash


# hadoop benchmark脚本文件（手动）
#
nrFiles=$1
size=$2

curdir=`pwd`
# command: ./hdfs-TestDFSIO.sh  write-append <nrFiles> <size>
#       size: 单位MB
  
#----------------------------TestDFSIO----------------------------#

# shellcheck disable=SC2046
start_time=$(date +%s)
echo "start time: " "$start_time"
mkdir testdfsio_log_"$nrFiles"_"$size"
hadoop jar "${HADOOP_HOME}"/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -write -nrFiles 1 -size $size"MB" -resFile "$curdir"/testdfsio_log_"$nrFiles"_"$size"/TestDFSIO-write-$nrFiles-1-$size.log
# shellcheck disable=SC2005
echo "write: $(date +%s)" >> time.log
for ((i=1; i<"$nrFiles"; i++))
do
  hadoop jar "${HADOOP_HOME}"/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -append -nrFiles 1 -size $size"MB" -resFile "$curdir"/testdfsio_log_"$nrFiles"_"$size"/TestDFSIO-append-$nrFiles-$i-$size.log
  echo "append $i: $(date +%s)" >> time.log
done
# shellcheck disable=SC2046
end_time=$(date +%s)
echo "HDFS start time: " "$start_time"
echo "HDFS end time:   " "$end_time"
# hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -clean

