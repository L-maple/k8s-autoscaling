#! /bin/bash

echo "STEP1 delete crd resources."
kubectl delete -f manifests/crd.yaml

echo "STEP2 delete tidb-admin namespace"
# helm uninstall tidb-operator 

kubectl delete namespace tidb-admin
