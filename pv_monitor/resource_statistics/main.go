package statistics

func main() {
	podStatistics := PodStatistics{
		Endpoint:  "http://localhost:9090",
		PodName:   "hdfs-datanode-0",
		Namespace: "monitoring",
	}

	podStatistics.GetLastDiskUtilization()
	podStatistics.GetLastMemoryUsage()
}
