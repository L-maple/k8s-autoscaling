#! /bin/bash

echo "STEP1 delete crd resources."
kubectl delete -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.1.11/manifests/crd.yaml

echo "STEP2 delete tidb-admin namespace"
kubectl delete namespace tidb-admin
