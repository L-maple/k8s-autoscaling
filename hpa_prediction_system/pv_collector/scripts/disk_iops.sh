#! /bin/bash

cmd="iostat -N | grep $1 | awk '{print \$2}'"
eval "$cmd"
