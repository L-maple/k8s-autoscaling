#! /bin/bash

target=$1
cmd="df --output=pcent,target | grep $target | awk '{print \$1}'"
eval "$cmd"