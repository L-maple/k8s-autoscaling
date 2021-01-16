# 删除可执行文件
echo "STEP 1: delete previous executable files."
rm -f pv_exporter

# 拉取新的主分支
echo "STEP 2: git pull new branch."
cd ../..
git pull
cd pv_monitor/pv_exporter/

# 构建新的可执行文件
echo "STEP 3: go build new exec file."
go build .

# 执行新的可执行文件
echo "STEP 4: exec new file."
./pv_exporter -statefulset=ubuntu