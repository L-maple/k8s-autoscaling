package main

import rs "github.com/k8s-autoscaling/pv_monitor/resource_statistics"

var (
	/* HPA Finite State*/
	FreeState    = 0
	StressState  = 1
	ScaleUpState = 2
)

func getHpaActivityState() int {
	stsMutex.RLock()
	defer stsMutex.RUnlock()

	if stsInfoGlobal.Initialized == false {
		return FreeState
	}
	podNameAndInfo := stsInfoGlobal.GetPodInfos()

	podCounter := len(podNameAndInfo)
	var cpuUtilizationSlice    []float64
	var memoryUtilizationSlice []int64
	var diskUtilizationSlice   []float64
	for podName, _ := range podNameAndInfo {
		podStatisticsObj := rs.PodStatistics{
			Endpoint:  prometheusUrl,
			PodName:   podName,
			Namespace: namespaceName,
		}

		cpuUtilizationSlice    = append(cpuUtilizationSlice, podStatisticsObj.GetLastCpuUtilizationQuery())
		memoryUtilizationSlice = append(memoryUtilizationSlice, podStatisticsObj.GetLastMemoryUtilizationQuery())
		diskUtilizationSlice   = append(diskUtilizationSlice, podStatisticsObj.GetLastDiskUtilizationQuery())
	}

	//avgCpuUtilization    := getAvgFloat64Utilization(cpuUtilizations)
	//avgMemoryUtilization := getAvgInt64Utilization(memoryUtilizations)
	avgDiskUtilization     := getAvgFloat64Utilization(diskUtilizationSlice)
	aboveNumber := getAboveUtilizationNumber(diskUtilizationSlice, 0.85)
	if podCounter - aboveNumber < 3 || avgDiskUtilization < 0.8 {
		return ScaleUpState
	}

	return FreeState
}

