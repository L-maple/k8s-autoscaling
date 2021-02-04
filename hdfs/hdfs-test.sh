#! /bin/bash

# 将benchmark相关文件复制到hdfs-client容器内部
kubectl -n monitoring cp ./benchmark/ hdfs-client:/benchmark/

# 登录容器查看下
kubectl exec -it -n monitoring hdfs-client -- /bin/bash

