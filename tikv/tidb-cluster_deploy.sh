#! /bin/bash

namespace="tidb-cluster"
cluster_name="advanced-tidb"

echo "STEP1 create the namespace tidb-cluster..."
kubectl create namespace $namespace

echo "STEP2 deploy the tidb cluster..."
kubectl apply -f tidb-cluster/tidb-cluster.yaml -n $namespace

echo "STEP3 check the status of pod..."
kubectl get po -n $namespace -l app.kubernetes.io/instance=$cluster_name


