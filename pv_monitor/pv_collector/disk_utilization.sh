#! /bin/bash

target=$1
result=`df --output=pcent,target | grep $target | awk '{print \$1}'`
echo $result