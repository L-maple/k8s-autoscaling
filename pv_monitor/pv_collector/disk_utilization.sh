#! /usr/bin sh

df --output=pcent,target | grep "$1" | awk '{print $1}'