# Build the docker image of HDFS
You can get the image from the DockerHub directly.
If you want to build by yourself, you can refer to below command. But you should download the package of Hadoop2.9.0 and JDK1.8 from their official website, and named the package as `hadoop-2.9.0.tar.gz` and `jdk-8u11-linux-x64.tar.gz`.
```shell
cd docker-build
docker build -t jinbozi/hdfs:0.2 .
```
> https://hub.docker.com/repository/docker/jinbozi/hdfs

# Deploy the HDFS
```shell
kubectl apply -f hpa-deploy/
kubectl get all -n monitoring
```

# Verify the HDFS cluster
You can visit the HDFS DashBoard with this address: http://localhost:32007/
You can enter into container by `kubectl exec -it -n monitoring pod/hdfs-datanode-0 -- /bin/bash` and use `hdfs dfs -ls /` to list HDFS files.

# Scale manual
```shell
kubectl scale sts hdfs-datanode --replicas=4
```

# Verift the HPA

Enter into container and put file to cluster.(The target value have been set as 3GiB, so we must fill at least 4GiB file.)
```
kubectl exec -it -n monitoring pod/hdfs-datanode-0 -- /bin/bash
dd if=/dev/zero of=test.file bs=1MiB count=4096
hdfs dfs -put test.file /
```

You can find the replicas of datanode will increase.
```
kubectl get statefulset.apps/hdfs-datanode -n monitoring
kubectl get hpa -n monitoring
```

Or visit the http://localhost:32007/

# Prometheus Export

You can use my custom Java prometheus export to get system directory size.([link](./docker-build/prometheus_export/))
> java -cp .:lib/* DirectorySizeExport [port] [path] [dir]
> port: the expose port to prometheus metrics, default value is 10318
> path: the path of du command, default value is /bin/du
> dir: the directory that you want to monitor its size, default value is /
> You must fill the previous argument if you want to custom the after argument.

