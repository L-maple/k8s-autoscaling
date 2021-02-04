#!/bin/bash


# hadoop benchmark脚本文件（自动）
#       
curdir=`pwd`
# command: ./mrbench.sh  <numRuns>  <maps> <reduces> <inputLines>
#       eg: ./mrbench.sh  10 10 5 10 
  
#----------------------------mrbench----------------------------#

mkdir mrbench_log
for a in {1..10}
    do
     for b in {10..100..10}
        do
        for c in {10..100..10}
            do
                for d in {10..100..10}
                    do
                    
    hadoop jar ${HADOOP_HOME}/share/hadoop/mapreduce/hadoop-mapreduce-client-jobclient-2.9.0-tests.jar mrbench -numRuns $a -maps $b -reduces $c -inputLines $d -inputType descending | tee ./mrbench_log/mrbench-$a-$b-$c-$d.log
   
                    done

            done

        done
    done
