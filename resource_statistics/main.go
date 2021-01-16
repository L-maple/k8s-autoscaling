package main

import (
	"flag"
)

var (
	endpoint        string   /* the prometheus http://url:port*/
	statefulSetName string
	namespace       string
	interval        int
)


func init() {
	flag.StringVar(&endpoint, "endpoint", "http://127.0.0.1:9090", "promethues url")
	flag.StringVar(&statefulSetName, "statefulset", "default", "statefulset Name")
	flag.StringVar(&namespace, "namespace", "default", "namespace")
	flag.IntVar(&interval, "interval", 5, "interval time")
}


func main() {
	podInfo := PodStatistics{}

	podInfo.GetAvgCpuUtilizationQuery()
}

