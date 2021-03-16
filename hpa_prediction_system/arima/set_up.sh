#! /bin/bash

# 构建镜像
docker build -t aliuchangjie/arima:latest .
echo "STEP1 docker build successfully."

# 登录docker
docker login
echo "STEP2 docker login successfully."

# 将镜像推送到仓库
docker push aliuchangjie/arima:latest
echo "STEP3 docker push successfully."

# 删除遗留的yaml部署资源
kubectl delete -f manifests/arima-deployment.yaml
kubectl delete -f manifests/arima-service.yaml
echo "STEP4 kubectl delete successfully."

# 部署yaml文件
kubectl apply -f manifests/arima-deployment.yaml
kubectl apply -f manifests/arima-service.yaml
echo "STEP5 kubectl apply successfully"

