package main

import rs "github.com/k8s-autoscaling/pv_monitor/resource_statistics"


func getHpaActivityState() int {
	stsMutex.RLock()
	defer stsMutex.RUnlock()

	if stsInfoGlobal.Initialized == false {
		return FREE_STATE
	}
	podNameAndInfo := stsInfoGlobal.GetPodInfos()

	podCounter := len(podNameAndInfo)
	var cpuUtilizations    []float64
	var memoryUtilizations []int64
	var diskUtilizations   []float64
	for podName, _ := range podNameAndInfo {
		podStatisticsObj := rs.PodStatistics{
			Endpoint:  prometheusUrl,
			PodName:   podName,
			Namespace: namespaceName,
		}

		cpuUtilizations    = append(cpuUtilizations, podStatisticsObj.GetLastCpuUtilizationQuery())
		memoryUtilizations = append(memoryUtilizations, podStatisticsObj.GetLastMemoryUtilizationQuery())
		diskUtilizations   = append(diskUtilizations, podStatisticsObj.GetLastDiskUtilizationQuery())
	}

	//avgCpuUtilization    := getAvgFloat64Utilization(cpuUtilizations)
	//avgMemoryUtilization := getAvgInt64Utilization(memoryUtilizations)
	avgDiskUtilization     := getAvgFloat64Utilization(diskUtilizations)
	aboveNumber := getAboveUtilizationNumber(diskUtilizations, 0.85)
	if podCounter - aboveNumber < 3 || avgDiskUtilization < 0.8 {
		return SCALEUP_STATE
	}

	return FREE_STATE
}

