#! /bin/bash

cmd="iostat -N | grep $1 | awk '{print \$4}'"
eval "$cmd"
