#! /bin/bash

echo "STEP1 create crd"
kubectl apply -f manifests/crd.yaml

echo "STEP2 download the operator"
mkdir -p ${HOME}/tidb-operator && \
helm inspect values pingcap/tidb-operator --version=v1.1.11 > ./tidb-operator/values-tidb-operator.yaml

echo "STEP3 configure the image"
sed -i 's?k8s.gcr.io?registry.cn-hangzhou.aliyuncs.com/google_containers?g' ./tidb-operator/values-tidb-operator.yaml

echo "STEP4 deploy the tidb-operator"
kubectl create namespace tidb-admin
helm install tidb-operator pingcap/tidb-operator --namespace=tidb-admin --version=v1.1.11 -f ./tidb-operator/values-tidb-operator.yaml 

sleep 5

kubectl get po -n tidb-admin -l app.kubernetes.io/name=tidb-operator

