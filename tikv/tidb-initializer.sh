#! /bin/bash

namespace="tidb-cluster"
root_password="root"

kubectl create secret generic tidb-secret --from-literal=root=$root_password --namespace=$namespace

kubectl apply -f tidb-cluster/tidb-initializer.yaml
