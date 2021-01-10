package main

import "strconv"

type PodStatistics struct {
	podName, namespace string
}

func (c PodStatistics) GetAvgCpuUtilizationQuery(durationSeconds int) string {
	query := "sum(rate(container_cpu_usage_seconds_total{" +
			 "image!=\"\", pod=" + c.podName + ", namespace=" + c.namespace +
			 "}[" + strconv.Itoa(durationSeconds) + "s]))"
	return query
}


