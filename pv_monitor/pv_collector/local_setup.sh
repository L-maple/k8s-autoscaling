#! /bin/bash

# 前提：
# 将main.go中 `flag.StringVar(&serverAddress, "serverAddress", "localhost:30002", "hpa-exporter comm address")` 取消注释
# 将main.go中 `flag.StringVar(&serverAddress, "serverAddress", "http://hpa-exporter-service.monitoring.svc:30002/", "hpa-exporter comm address")` 注释掉

git pull
echo "STEP1. git pull SUCCESS"

rm -f pv_collector
echo "STEP2. rm the pv_collector SUCCESS"

go build .
echo "STEP3. go build . SUCCESS"

./pv_collector
echo "STEP4. ./pv_collector SUCCESS"

