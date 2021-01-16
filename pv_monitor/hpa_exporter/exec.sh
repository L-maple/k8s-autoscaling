# 删除可执行文件
echo "STEP 1: delete previous executable files."
rm -f hpa_exporter

# 拉取新的主分支
echo "STEP 2: git pull new branch."
cd ../..
git pull
cd pv_monitor/hpa_exporter/

# 构建新的可执行文件
echo "STEP 3: go build new exec file."
go build .

# 执行新的可执行文件
echo "STEP 4: exec new file."
./hpa_exporter -statefulset=ubuntu

# 生成新的镜像，并打包上传

