package main

import (
	"fmt"
	rs "github.com/k8s-autoscaling/pv_monitor/resource_statistics"
	"time"
)

var (
	/* 副本数量 */
	ReplicasAmount = 3

	/* HPA Finite State*/
	FreeState      = 0
	StressState    = 1
	ScaleUpState   = 2
)

func getHpaActivityState() int {
	stsMutex.RLock()
	defer stsMutex.RUnlock()
	// 如果 stsInfoGlobal还没初始化，那么直接返回 FreeState
	if stsInfoGlobal.Initialized == false {
		fmt.Println("stsInfoGlobal.Initialized: ", stsInfoGlobal.Initialized)
		printStatefulSetState(stsInfoGlobal)
		return FreeState
	}
	podNameAndInfo  := stsInfoGlobal.GetPodInfos()
	cpumilliLimit   := stsInfoGlobal.GetCpuMilliLimit()
	memoryByteLimit := stsInfoGlobal.GetMemoryByteLimit()

	podCounter := len(podNameAndInfo)
	var cpuUsageSlice          []float64
	var memoryUsageSlice       []int64
	var diskUtilizationSlice   []float64
	for podName, _ := range podNameAndInfo {
		podStatisticsObj := rs.PodStatistics{
			Endpoint:  prometheusUrl,
			PodName:   podName,
			Namespace: namespaceName,
		}

		cpuUsageSlice        = append(cpuUsageSlice, podStatisticsObj.GetLastCpuUsage())
		memoryUsageSlice     = append(memoryUsageSlice, podStatisticsObj.GetLastMemoryUsage())
		diskUtilizationSlice = append(diskUtilizationSlice, podStatisticsObj.GetLastDiskUtilization())
	}

	// TODO: 设置稳定窗口计时器
	// 得到CPU的平均使用量
	avgCpuUsage    := getAvgFloat64(cpuUsageSlice)
	avgCpuUtilization := avgCpuUsage / float64(cpumilliLimit)

	// 得到Memory的平均使用量
	avgMemoryUsage := getAvgInt64(memoryUsageSlice)
	avgMemoryUtilization := float64(avgMemoryUsage) / float64(memoryByteLimit)

	avgDiskUtilization    := getAvgFloat64(diskUtilizationSlice)
	aboveNumber := getGreaterThanStone(diskUtilizationSlice, 0.85)

	printCurrentState(avgCpuUtilization, avgMemoryUtilization, avgDiskUtilization, podCounter, aboveNumber)

	if podCounter - aboveNumber < ReplicasAmount || avgDiskUtilization > 0.8 {
		return ScaleUpState
	}

	return FreeState
}

func printCurrentState(avgCpuUtilization, avgMemoryUtilization, avgDiskUtilization float64,
						podCounter, aboveNumber int) {
	fmt.Printf("++++++++++++++++++++++++++++++++++++\n")
	fmt.Printf("[INFO] %v\n", time.Now())
	stsMutex.RLock()
	stsTmp := stsInfoGlobal
	stsMutex.RUnlock()
	printStatefulSetState(stsTmp)

	fmt.Printf("avgCpuUtilization: %-30.3f, avgMemoryUtilization: %-30.3f, avgDiskUtilization: %-30.3f\n",
					avgCpuUtilization, avgMemoryUtilization, avgDiskUtilization)
	fmt.Printf("pod Numbers: %d, aboveNumber: %d\n", podCounter, aboveNumber)
	fmt.Printf("====================================\n\n")
}

func printStatefulSetState(stsInfo StatefulSetInfo) {
	fmt.Printf("%-40s, %-40s, %-40s\n", "PodName", "PvcName", "PvName")
	for podName, podInfo := range stsInfo.GetPodInfos() {
		fmt.Printf("%-40s ", podName)

		for _, pvcName := range podInfo.PVCNames {
			fmt.Printf("%-40s ", pvcName)
		}
		fmt.Printf("; ")

		for _, pvName := range podInfo.PVNames {
			fmt.Printf("%-40s ", pvName)
		}
		fmt.Printf("; \n")
	}
}