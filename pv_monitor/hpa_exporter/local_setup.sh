#! /bin/bash

# 前提：
# 将main.go中 `clientSet := getInClusterClientSet()` 注释掉
# 将main.go中 `clientSet := getClientSet()` 取消注释
# 将main.go中 `flag.StringVar(&prometheusUrl, "prometheus-url", "http://prometheus-k8s.monitoring.svc:9090/", "promethues url")`注释掉
# 将main.go中 `flag.StringVar(&prometheusUrl, "prometheus-url", "http://127.0.0.1:9090/", "promethues url")`取消注释

rm -f hpa_exporter
echo "STEP1. rm the hpa_exporter SUCCESS"

git pull
echo "STEP2. git pull SUCCESS"

go build .
echo "STEP3. go build . SUCCESS"

./hpa_exporter -namespace=monitoring -statefulset=hdfs-datanode
echo "STEP4. execute SUCCESS"
