package main

import (
	"fmt"
	rs "github.com/k8s-autoscaling/pv_monitor/resource_statistics"
)

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

	// TODO: 设置稳定窗口
	// 得到CPU的使用比率
	avgCpuUtilization    := getAvgFloat64(cpuUtilizationSlice)

	// 得到Memory的使用比率
	avgMemoryUtilization := getAvgInt64(memoryUtilizationSlice)

	fmt.Println("cpuUtilization: ", avgCpuUtilization)
	fmt.Println("memoryUtilization: ", avgMemoryUtilization)

	avgDiskUtilization    := getAvgFloat64(diskUtilizationSlice)
	aboveNumber := getGreaterThanStone(diskUtilizationSlice, 0.85)
	if podCounter - aboveNumber < 3 || avgDiskUtilization < 0.8 {
		return ScaleUpState
	}

	return FreeState
}

