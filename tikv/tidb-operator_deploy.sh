#! /bin/bash

chart_version="v1.1.11"

echo "STEP1 create CRD"
kubectl apply -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.1.11/manifests/crd.yaml

echo "STEP2 download the values-tidb-operator.yaml."
mkdir -p tidb-operator && \
helm inspect values pingcap/tidb-operator --version=${chart_version} > tidb-operator/values-tidb-operator.yaml

echo "STEP3 configure the image source"
sed -i 's?k8s.gcr.io?registry.cn-hangzhou.aliyuncs.com/google_containers?g' tidb-operator/values-tidb-operator.yaml

echo "STEP4 deploy the tidb-operator"
kubectl create namespace tidb-admin
helm install tidb-operator pingcap/tidb-operator --namespace=tidb-admin --version=${chart_version} -f tidb-operator/values-tidb-operator.yaml 

sleep 5

kubectl get po -n tidb-admin -l app.kubernetes.io/name=tidb-operator

