#! /bin/bash

git pull
echo "STEP1. git pull SUCCESS"

rm -f pv_collector
echo "STEP2. rm the pv_collector SUCCESS"

go build .
echo "STEP3. go build . SUCCESS"

./pv_collector
echo "STEP4. ./pv_collector SUCCESS"

