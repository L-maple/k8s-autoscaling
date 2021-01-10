#! /bin/bash

target=$1
echo $target
cmd="df --output=pcent,target | grep $target | awk '{print \$1}'"
eval "$cmd"