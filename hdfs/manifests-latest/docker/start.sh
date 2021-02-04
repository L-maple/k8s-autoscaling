#!/usr/bin/env bash

function setDefault(){
    if [[ "`eval echo '$'"$1"`" = "" ]]; then
        export `eval echo "$1"`=$2
    fi
}

setDefault hadoop_core_fs_defaultFS "hdfs://0.0.0.0:9000"
setDefault hadoop_core_dfs_namenode_rpc_bind_host "0.0.0.0"
setDefault hadoop_core_hadoop_tmp_dir "/tmp/hadoop"

setDefault hadoop_hdfs_dfs_namenode_name_dir "file:///custom/hdfs_dir/namenode"
setDefault hadoop_hdfs_dfs_datanode_data_dir "file:///custom/hdfs_dir/namenode"
setDefault hadoop_hdfs_dfs_replication "3"
setDefault hadoop_hdfs_dfs_datanode_balance_bandwidthPerSec "10g"

setDefault hadoop_yarn_yarn_resourcemanager_hostname "0.0.0.0"
setDefault hadoop_yarn_yarn_application_classpath $(hadoop classpath)
setDefault hadoop_yarn_yarn_scheduler_minimum_allocation_mb "1024"
setDefault hadoop_yarn_yarn_scheduler_maximum_allocation_mb "2048"
setDefault hadoop_yarn_yarn_scheduler_minimum_allocation_vcores "1"
setDefault hadoop_yarn_yarn_scheduler_maximum_allocation_vcores "4"
setDefault hadoop_yarn_yarn_nodemanager_resource_memory_mb "2048"

setDefault hadoop_mapred_yarn_app_mapreduce_am_env "HADOOP_MAPRED_HOME=/custom/hadoop-2.9.0"
setDefault hadoop_mapred_mapreduce_map_env "HADOOP_MAPRED_HOME=/custom/hadoop-2.9.0"
setDefault hadoop_mapred_mapreduce_reduce_env "HADOOP_MAPRED_HOME=/custom/hadoop-2.9.0"
setDefault hadoop_mapred_mapreduce_application_classpath $(hadoop classpath)
setDefault hadoop_mapred_mapreduce_map_memory_mb "1024"
setDefault hadoop_mapred_mapreduce_reduce_memory_mb "1024"
setDefault hadoop_mapred_mapred_child_java_opts "-Xmx200m"
setDefault hadoop_mapred_mapreduce_map_cpu_vcores "1"
setDefault hadoop_mapred_mapreduce_reduce_cpu_vcores "1"
setDefault hadoop_mapred_mapreduce_task_io_sort_mb "100"
setDefault hadoop_mapred_mapreduce_map_sort_spill_percent "0.8"
setDefault hadoop_mapred_mapreduce_map_maxattempts "4"
setDefault hadoop_mapred_mapreduce_reduce_maxattempts "4"
setDefault hadoop_mapred_mapreduce_job_maxtaskfailures_per_tracker "0"
setDefault hadoop_mapred_mapreduce_task_timeout "600000"

setDefault hadoop_custom_balance_time_interval "120s"
setDefault hadoop_custom_balance_threshold "15"

sed -i "s#hadoop_core_fs_defaultFS#${hadoop_core_fs_defaultFS}#g" $HADOOP_HOME/etc/hadoop/core-site.xml
sed -i "s#hadoop_core_dfs_namenode_rpc-bind-host#${hadoop_core_dfs_namenode_rpc_bind_host}#g" $HADOOP_HOME/etc/hadoop/core-site.xml
sed -i "s#hadoop_core_hadoop_tmp_dir#${hadoop_core_hadoop_tmp_dir}#g" $HADOOP_HOME/etc/hadoop/core-site.xml

sed -i "s#hadoop_hdfs_dfs_namenode_name_dir#${hadoop_hdfs_dfs_namenode_name_dir}#g" $HADOOP_HOME/etc/hadoop/hdfs-site.xml
sed -i "s#hadoop_hdfs_dfs_datanode_data_dir#${hadoop_hdfs_dfs_datanode_data_dir}#g" $HADOOP_HOME/etc/hadoop/hdfs-site.xml
sed -i "s#hadoop_hdfs_dfs_replication#${hadoop_hdfs_dfs_replication}#g" $HADOOP_HOME/etc/hadoop/hdfs-site.xml
sed -i "s#hadoop_hdfs_dfs_datanode_balance_bandwidthPerSec#${hadoop_hdfs_dfs_datanode_balance_bandwidthPerSec}#g" $HADOOP_HOME/etc/hadoop/hdfs-site.xml

sed -i "s#hadoop_yarn_yarn_resourcemanager_hostname#${hadoop_yarn_yarn_resourcemanager_hostname}#g" $HADOOP_HOME/etc/hadoop/yarn-site.xml
sed -i "s#hadoop_yarn_yarn_application_classpath#${hadoop_yarn_yarn_application_classpath}#g" $HADOOP_HOME/etc/hadoop/yarn-site.xml
sed -i "s#hadoop_yarn_yarn_scheduler_minimum-allocation-mb#${hadoop_yarn_yarn_scheduler_minimum_allocation_mb}#g" $HADOOP_HOME/etc/hadoop/yarn-site.xml
sed -i "s#hadoop_yarn_yarn_scheduler_maximum-allocation-mb#${hadoop_yarn_yarn_scheduler_maximum_allocation_mb}#g" $HADOOP_HOME/etc/hadoop/yarn-site.xml
sed -i "s#hadoop_yarn_yarn_scheduler_minimum-allocation-vcores#${hadoop_yarn_yarn_scheduler_minimum_allocation_vcores}#g" $HADOOP_HOME/etc/hadoop/yarn-site.xml
sed -i "s#hadoop_yarn_yarn_scheduler_maximum-allocation-vcores#${hadoop_yarn_yarn_scheduler_maximum_allocation_vcores}#g" $HADOOP_HOME/etc/hadoop/yarn-site.xml
sed -i "s#hadoop_yarn_yarn_nodemanager_resource_memory-mb#${hadoop_yarn_yarn_nodemanager_resource_memory_mb}#g" $HADOOP_HOME/etc/hadoop/yarn-site.xml

sed -i "s#hadoop_mapred_yarn_app_mapreduce_am_env#${hadoop_mapred_yarn_app_mapreduce_am_env}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_map_env#${hadoop_mapred_mapreduce_map_env}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_reduce_env#${hadoop_mapred_mapreduce_reduce_env}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_application_classpath#${hadoop_mapred_mapreduce_application_classpath}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_map_memory_mb#${hadoop_mapred_mapreduce_map_memory_mb}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_reduce_memory_mb#${hadoop_mapred_mapreduce_reduce_memory_mb}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapred_child_java_opts#${hadoop_mapred_mapred_child_java_opts}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_map_cpu_vcores#${hadoop_mapred_mapreduce_map_cpu_vcores}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_reduce_cpu_vcores#${hadoop_mapred_mapreduce_reduce_cpu_vcores}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_task_io_sort_mb#${hadoop_mapred_mapreduce_task_io_sort_mb}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_map_sort_spill_percent#${hadoop_mapred_mapreduce_map_sort_spill_percent}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_map_maxattempts#${hadoop_mapred_mapreduce_map_maxattempts}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_reduce_maxattempts#${hadoop_mapred_mapreduce_reduce_maxattempts}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_job_maxtaskfailures_per_tracker#${hadoop_mapred_mapreduce_job_maxtaskfailures_per_tracker}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml
sed -i "s#hadoop_mapred_mapreduce_task_timeout#${hadoop_mapred_mapreduce_task_timeout}#g" $HADOOP_HOME/etc/hadoop/mapred-site.xml

tmp=$image_function
functions=(${tmp//|/ })  
for function in ${functions[@]}
do
    echo $function
    if [[ "$function" = "namenode" ]]; then
        echo "Start namenode ..."
        if [ ! -d "/custom/hdfs_dir/namenode" ]
        then
            echo "Format namenode ..."
            hdfs namenode -format
        fi
        hadoop-daemon.sh start namenode
    elif [[ "$function" = "datanode" ]]; then
        echo "Start datanode ..."
        hadoop-daemon.sh start datanode
    elif [[ "$function" = "resourcemanager" ]]; then
        echo "Start resourcemanager ..."
        yarn-daemon.sh start resourcemanager
    elif [[ "$function" = "nodemanager" ]]; then
        echo "Start nodemanager ..."
        yarn-daemon.sh start nodemanager
    elif [[ "$function" = "prometheus" ]]; then
        echo "Export prometheus metrics ..."
        java -cp /custom/prometheus_export:/custom/prometheus_export/lib/* HDFSDiskExport 10318 /custom/prometheus_export/get_system_metrics.sh /custom/hdfs_dir &
    elif [[ "$function" = "execute" ]]; then
        echo "Execute $command ..."
        $command
    elif [[ "$function" = "ssh" ]]; then
        echo "Start ssh service ..."
        /etc/init.d/ssh start
    elif [[ "$function" = "balance" ]]; then
        echo "Balance the hdfs ..."
        while true
        do
            sleep $hadoop_custom_balance_time_interval
            hdfs balancer -threshold $hadoop_custom_balance_threshold
        done
    fi
done 

while true; do sleep 10000; done