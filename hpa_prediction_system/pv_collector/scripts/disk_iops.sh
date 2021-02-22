#! /bin/bash

#cmd="iostat -N | grep $1 | awk '{print \$2}'"
#eval "$cmd"

device_name="lvdisplay -c | grep $1 |  awk -F ':' '{print 'dm-'\$13}'"
cmd="iostat | grep $device_name | awk '{print \$2}'"
eval "$cmd"