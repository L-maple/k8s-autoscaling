#!/bin/bash


# hadoop benchmark脚本文件（手动）
#       
a=$1
b=$2
c=$3
d=$4
curdir=`pwd`
# command: ./mrbench.sh  <numRuns>  <maps> <reduces> <inputLines>
#       eg: ./mrbench.sh  10 10 5 10 
  
#----------------------------mrbench----------------------------#

#hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar TestDFSIO -$a -nrFiles $b -size $c"MB" -resFile $curdir/TestDFSIO-$a-$b-$c.log
hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar mrbench -numRuns $a -maps $b -reduces $c -inputLines $d -inputType descending | tee mrbench-$a-$b-$c-$d.log