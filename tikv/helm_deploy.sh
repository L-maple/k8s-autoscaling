#! /bin/bash

echo "STEP1: download the source."
wget https://get.helm.sh/helm-v3.4.1-linux-amd64.tar.gz
rm -rf /usr/local/helm
mkdir /usr/local/helm
tar zxvf  helm-v3.4.1-linux-amd64.tar.gz -C /usr/local/helm
rm -f helm-v3.4.1-linux-amd64.tar.gz

echo "STEP2: copy the binary to /usr/local/bin."
rm /usr/local/bin/helm -f
cp /usr/local/helm/linux-amd64/helm /usr/local/bin

echo "STEP3: check the helm version."
helm version

echo "STEP4: configure helm repo."
helm repo add pingcap https://charts.pingcap.org/
helm search repo pingcap

echo "STEP5: update local cache."
helm repo update


