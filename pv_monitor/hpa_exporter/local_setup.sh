#! /bin/bash

rm -f hpa_exporter
echo "STEP1. rm the hpa_exporter SUCCESS"

git pull
echo "STEP2. git pull SUCCESS"

go build .
echo "STEP3. go build . SUCCESS"

./hpa_exporter -namespace=monitoring -statefulset=hdfs-datanode
echo "STEP4. execute SUCCESS"
