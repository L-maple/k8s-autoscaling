#!/bin/bash


# hadoop benchmark脚本文件自动）
#       
a=$1

curdir=`pwd`
# command: ./Terasort.sh  
#       eg: ./Terasort.h
  
#数据生成目录在/benchmarks/test_data
#结果生成目录在/benchmarks/terasort-output
#未计算数据生成时间
#
#----------------------------Terasort----------------------------#

mkdir Terasort_log
for a in {10..1000..50}
    do
    hadoop fs -rmr /benchmarks/test_data
    hadoop fs -rmr /benchmarks/terasort-output

    hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-examples-2.9.0.jar teragen $a /benchmarks/test_data

    timer_start=`date "+%Y-%m-%d %H:%M:%S"`

    hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-examples-2.9.0.jar terasort /benchmarks/test_data /benchmarks/terasort-output

    timer_end=`date "+%Y-%m-%d %H:%M:%S"`

    duration=`echo $(($(date +%s -d "${timer_end}") - $(date +%s -d "${timer_start}"))) | awk '{t=split("60 s 60 m 24 h 999 d",a);for(n=1;n<t;n+=2){if($1==0)break;s=$1%a[n]a[n+1]s;$1=int($1/a[n])}print s}'`
    echo "耗时： $duration" >./Terasort_log/Terasort_$a.log

    done
