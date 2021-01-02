#! /bin/sh

# 删除遗留的可执行文件
rm -rf pv_exporter
echo "Step1: rm SUCCESS."

# 编译生成可执行文件
go build -o pv_exporter .
echo "Step2: go build SUCCESS."

# 构建镜像
docker build -t aliuchangjie/pv_exporter .
echo "Step3: docker build SUCCESS."

# 登录
docker login
echo "Step4: docker login SUCCESS."

# 将镜像推送到仓库
docker push aliuchangjie/pv_exporter
echo "Step5: docker push SUCCESS."

# 删除遗留的yaml部署资源
kubectl delete -f pv_exporter.yaml
echo "Step6: kubectl delete SUCCESS."

# 部署yaml文件
kubectl apply -f pv_exporter.yaml
echo "Step7: kubectl apply SUCCESS."
