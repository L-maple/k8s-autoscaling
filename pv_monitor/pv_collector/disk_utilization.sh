#! /bin/bash

result=`df --output=pcent,target | grep "$1" | awk '{print $1}'`
echo $result