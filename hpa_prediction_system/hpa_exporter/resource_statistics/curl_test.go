package statistics

import (
	"fmt"
	"testing"
)

func TestPromCurl_Get(t *testing.T) {
	podStatistics := PodStatistics{
		Endpoint:  "http://localhost:9090",
		PodName:   "hdfs-datanode-0",
		Namespace: "monitoring",
	}

	fmt.Println(podStatistics.GetLastDiskUtilization())
	fmt.Println(podStatistics.GetLastMemoryUsage())
	fmt.Println(podStatistics.GetLastCpuUtilization())
}
