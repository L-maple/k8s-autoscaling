#!/bin/bash

echo "Step1: 删除之前不必要的yaml资源"
kubectl delete -f manifests-latest/
sleep 10

echo "Step2: 删除所有的pvc"
kubectl delete -n monitoring pvc --all

echo "Step3: 部署新的yaml资源"
kubectl apply -f manifests-latest/

echo "Step4: 记得删除worker节点上的无用pv\n"

echo "Step5: 等待，不然hpa可能会做出错误的扩容决定"
sleep 300

# 将benchmark相关文件复制到hdfs-client容器内部
kubectl -n monitoring cp ./benchmark/ hdfs-client:/benchmark/

# 登录容器查看下
kubectl exec -it -n monitoring hdfs-client -- /bin/bash



