#! /bin/sh

# 前提：
# 将main.go中 `flag.StringVar(&serverAddress, "serverAddress", "localhost:30002", "hpa-exporter comm address")` 注释掉
# 将main.go中 `flag.StringVar(&serverAddress, "serverAddress", "http://hpa-exporter-service.monitoring.svc:30002/", "hpa-exporter comm address")` 取消注释

# Git pull
git pull

# 删除遗留的可执行文件
rm -f pv_collector
echo "Step1: rm SUCCESS."

# 删除遗留的yaml部署资源
kubectl delete -f manifests/pv_collector.yaml
echo "Step2: kubectl delete SUCCESS."

# 编译生成可执行文件
go build -o pv_collector .
echo "Step3: go build SUCCESS."

# 构建镜像
docker rmi aliuchangjie/pv_collector:latest
docker build -t aliuchangjie/pv_collector:latest .
echo "Step4: docker build SUCCESS."

# 登录
docker login
echo "Step5: docker login SUCCESS."

# 将镜像推送到仓库
docker push aliuchangjie/pv_collector:latest
echo "Step6: docker push SUCCESS."

# 部署yaml文件
kubectl apply -f manifests/pv_collector.yaml
echo "Step7: kubectl apply SUCCESS."
