#! /bin/bash

echo "STEP1. rm the hpa_exporter"
rm -f hpa_exporter

echo "STEP2. git pull"
git pull

echo "STEP3. go build ."
go build .

echo "STEP4. execute"
./hpa_exporter -namespace=monitoring -statefulset=hdfs-datanode
