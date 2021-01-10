#! /bin/bash

cmd="iostat -N | grep $1 | awk '{print \$3}'"
eval "$cmd"
