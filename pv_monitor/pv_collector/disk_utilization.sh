#! /bin/bash

echo "df --output=pcent,target | grep $1 | awk '{print $1}'"