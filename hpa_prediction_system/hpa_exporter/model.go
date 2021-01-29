package main

import (
	"fmt"
	rs "github.com/k8s-autoscaling/hpa_prediction_system/resource_statistics"
	"log"
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

	fmt.Println("stsInfoGlobal.Initialized: ", stsInfoGlobal.Initialized)
	// 如果 stsInfoGlobal还没初始化，那么直接返回 FreeState
	if stsInfoGlobal.Initialized == false {
		printStatefulSetState(&stsInfoGlobal)
		return FreeState
	}
	podNameAndInfo  := stsInfoGlobal.GetPodInfos()
	memoryByteLimit := stsInfoGlobal.GetMemoryByteLimit()

	podCounter := len(podNameAndInfo)
	var cpuUtilizationSlice    []float64
	var memoryUsageSlice       []int64
	var diskUtilizationSlice   []float64
	for podName, _ := range podNameAndInfo {
		podStatisticsObj := rs.PodStatistics{
			Endpoint:  prometheusUrl,
			PodName:   podName,
			Namespace: namespaceName,
		}

		cpuUtilizationSlice  = append(cpuUtilizationSlice, podStatisticsObj.GetLastCpuUtilization())
		memoryUsageSlice     = append(memoryUsageSlice, podStatisticsObj.GetLastMemoryUsage())
		diskUtilizationSlice = append(diskUtilizationSlice, podStatisticsObj.GetLastDiskUtilization())
	}

	// TODO: 设置稳定窗口计时器
	// 得到CPU的平均使用量
	avgCpuUtilization    := getAvgFloat64(cpuUtilizationSlice)

	// 得到Memory的平均使用量
	avgMemoryUsage := getAvgInt64(memoryUsageSlice)
	avgMemoryUtilization := float64(avgMemoryUsage) / float64(memoryByteLimit)

	avgDiskUtilization    := getAvgFloat64(diskUtilizationSlice)
	aboveCeilingNumber := getGreaterThanStone(diskUtilizationSlice, 0.85)

	printCurrentState(avgCpuUtilization, avgMemoryUtilization, avgDiskUtilization, podCounter, aboveCeilingNumber)

	if podCounter - aboveCeilingNumber < ReplicasAmount || avgDiskUtilization > 0.8 {
		return ScaleUpState
	}

	return FreeState
}

func printCurrentState(avgCpuUtilization, avgMemoryUtilization, avgDiskUtilization float64,
						podCounter, aboveCeilingNumber int) {
	fmt.Printf("++++++++++++++++++++++++++++++++++++\n")
	fmt.Printf("[INFO] %v\n", time.Now())

	printStatefulSetState(&stsInfoGlobal)

	fmt.Printf("avgCpuUtilization: %-30.6f, avgMemoryUtilization: %-30.6f, avgDiskUtilization: %-30.6f\n",
					avgCpuUtilization, avgMemoryUtilization, avgDiskUtilization)
	fmt.Printf("pod Numbers: %d, aboveCeilingNumber: %d\n", podCounter, aboveCeilingNumber)
	fmt.Printf("====================================\n\n")
}

func printStatefulSetState(stsInfo *StatefulSetInfo) {
	fmt.Printf("%-40s %-40s %-40s\n", "PodName", "PvcName", "PvName")
	for podName, podInfo := range stsInfo.GetPodInfos() {
		fmt.Printf("%-40s ", podName)

		for _, pvcName := range podInfo.PVCNames {
			fmt.Printf("%-40s ", pvcName)
		}

		for _, pvName := range podInfo.PVNames {
			fmt.Printf("%-40s ", pvName)
		}

		for index, pvName := range podInfo.PVNames {
			fmt.Printf("PV Name{%d}: %s: \n", index, pvName)
			diskIOPS, err        := getLastDiskIOPS(pvName)
			if err != nil {
				log.Fatal("getLastDiskIOPS: ", err)
			}
			diskReadMBPS, err    := getLastDiskReadMBPS(pvName)
			if err != nil {
				log.Fatal("getLastDiskReadMBPS: ", err)
			}
			diskWriteMBPS, err   := getLastWriteMBPS(pvName)
			if err != nil {
				log.Fatal("getLastWriteMBPS: ", err)
			}
			diskUtilization, err := getLastDiskUtilization(pvName)
			if err != nil {
				log.Fatal("getLastDiskUtilization: ", err)
			}
			fmt.Printf("diskIOPS: %-10.6f, diskReadMBPS: %-10.6f, diskWriteMBPS: %-10.6f, diskUtilization: %-10.6f\n\n",
				diskIOPS, diskReadMBPS, diskWriteMBPS, diskUtilization)
		}
	}
}