#!/bin/bash
set -u

du -s $1 | awk '{print $1}'

tmp=`df -l | grep $1 | awk '{print $5}'`
tmp=${tmp%?}
echo $tmp | awk '{print $1/100.0}'

tmp=`df -lh | grep $1 | awk -F'[ /]' '{print $4}'`
res=`iostat -N -m | grep $tmp`
echo $res | awk '{print $2}'
echo $res | awk '{print $3}'
echo $res | awk '{print $4}'
echo $res | awk '{print $5}'
echo $res | awk '{print $6}'