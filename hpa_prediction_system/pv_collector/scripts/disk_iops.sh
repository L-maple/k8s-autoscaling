#! /bin/bash

#cmd="iostat -N | grep $1 | awk '{print \$2}'"
#eval "$cmd"

device_name_command="lvdisplay -c | grep $1 |  awk -F ':' '{print 'dm-'\$13}'"
dm_name=$(eval "$device_name_command")
cmd="iostat | grep $dm_name | awk '{print \$2}'"
eval "$cmd"