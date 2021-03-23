#! /bin/sh

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
kubectl delete -f tidb-manifests/tidb-hpa-exporter.yaml
kubectl delete -f tidb-manifests/tidb-custom-hpa.yaml
# 删除可执行文件
rm -f hpa_exporter

echo "Step6: kubectl delete SUCCESS."

# 部署yaml文件
kubectl apply -f tidb-manifests/tidb-hpa-exporter.yaml
kubectl apply -f tidb-manifests/tidb-custom-hpa.yaml
echo "Step7: kubectl apply SUCCESS."

