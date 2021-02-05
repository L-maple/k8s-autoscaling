#!/bin/bash


# hadoop benchmark脚本文件（自动）
#       

# command: ./hdfs-TestDFSIO.sh  
#           <write/read>  write read
#           <nrFiles>   1-10
#           <size>      100-1100 step 100
#       
  
#----------------------------TestDFSIO----------------------------#

mkdir TestDFSIO_log
for i in {1..11}
do
   for j in {100..1100..100}
    do
      hadoop fs -rmr /benchmarks/TestDFSIO
      hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -write -nrFiles $i -size $j"MB" -resFile ./TestDFSIO_log/TestDFSIO-write-$i-$j.log
    
      hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -read -nrFiles $i -size $j"MB" -resFile ./TestDFSIO_log/TestDFSIO-read-$i-$j.log

    done
done




