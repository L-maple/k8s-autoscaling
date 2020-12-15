# Deploy the HDFS
```shell
kubectl apply -f deploy/
```

# Visit the web of NameNode
http://k8s-master-ip:32007/

# Scale manual
```shell
kubectl scale sts nwpu-datanode --replicas=4
```

# The change of HDFS cluster
befor scale : 3 pods and 3 DataNodes
![befor system](./pictures/prev1.png)
![befor web](./pictures/prev2.png)
after scale : 4 pods and 4 DataNodes
![after system](./pictures/after1.png)
![after web](./pictures/after2.png)
