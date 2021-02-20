#! /bin/sh

# 前提：
# 将main.go中 `clientSet := getInClusterClientSet()` 取消注释
# 将main.go中 `clientSet := getClientSet()` 注释掉
# 将main.go中 `flag.StringVar(&prometheusUrl, "prometheus-url", "http://prometheus-k8s.monitoring.svc:9090/", "promethues url")`取消注释
# 将main.go中 `flag.StringVar(&prometheusUrl, "prometheus-url", "http://127.0.0.1:9090/", "promethues url")` 注释掉

# Git pull
git pull

# 删除遗留的可执行文件
rm -rf hpa_exporter
echo "Step1: rm SUCCESS."

# 编译生成可执行文件
go build -o hpa_exporter .
echo "Step2: go build SUCCESS."

# 构建镜像
docker build -t aliuchangjie/hpa_exporter:latest .
echo "Step3: docker build SUCCESS."

# 登录
docker login
echo "Step4: docker login SUCCESS."

# 将镜像推送到仓库
docker push aliuchangjie/hpa_exporter:latest
echo "Step5: docker push SUCCESS."

# 删除遗留的yaml部署资源
kubectl delete -f manifests/hdfs-hpa-exporter.yaml
kubectl delete -f manifests/hdfs-custom-hpa.yaml
# 删除可执行文件
rm -f hpa_exporter

echo "Step6: kubectl delete SUCCESS."

# 部署yaml文件
kubectl apply -f manifests/hdfs-hpa-exporter.yaml
kubectl apply -f manifests/hdfs-custom-hpa.yaml
echo "Step7: kubectl apply SUCCESS."

