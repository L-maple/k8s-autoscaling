#! /bin/bash

namespace=tidb-cluster
clustername=advanced-tidb

mysql_port=`kubectl -n $namespace get svc ${clustername}-tidb -ojsonpath="{.spec.ports[?(@.name=='mysql-client')].nodePort}{'\n'}"`

mysql --host 127.0.0.1 --port $mysql_port -p"root" -u "root" 
