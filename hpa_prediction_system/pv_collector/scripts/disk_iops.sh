#! /bin/bash

#cmd="iostat -N | grep $1 | awk '{print \$2}'"
#eval "$cmd"

# shellcheck disable=SC2140
device_name_command="lvdisplay -c | grep $1 |  awk -F ':' '{print \$13}'"
dm_name="dm-"$(eval "$device_name_command")
#echo "$dm_name"
cmd="iostat | grep '$dm_name' | awk '{print \$2}'"
eval "$cmd"