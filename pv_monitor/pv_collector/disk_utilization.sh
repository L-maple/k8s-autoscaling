#! /bin/bash

cmd="df --output=pcent,target | grep $1 | awk '{print \$1 * 100}'"
eval "$cmd"