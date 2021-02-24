#! /bin/bash

# shellcheck disable=SC2140
device_name_command="lvdisplay -c | grep $1 |  awk -F ':' '{print \$13}'"
dm_name="dm-"$(eval "$device_name_command")
if [ "$dm_name" == "dm-" ]
then
  echo ""
else
  cmd="iostat | grep '$dm_name' | awk '{print \$2}'"
  eval "$cmd"
fi