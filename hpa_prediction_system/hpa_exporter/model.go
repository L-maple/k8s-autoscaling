package main

import (
	"fmt"
	rs "github.com/k8s-autoscaling/hpa_prediction_system/hpa_exporter/resource_statistics"
	"log"
	"time"
)

var (
	/* 副本数量 */
	ReplicasAmount = 3
)

func getHpaActivityState() int {
	// 如果 stsInfoGlobal还没初始化，那么直接返回 FreeState
	if stsInfoGlobal.Initialized == false {
		printStatefulSetState(stsInfoGlobal)

		return hpaFSM.GetState()
	}

	printCurrentState()

	return hpaFSM.GetState()
}

func printCurrentState() {
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

	// 得到CPU的平均使用率
	avgCpuUtilization    := getAvgFloat64(cpuUtilizationSlice)

	// 得到Memory的平均使用率
	avgMemoryUsage := getAvgInt64(memoryUsageSlice)
	avgMemoryUtilization := float64(avgMemoryUsage) / float64(memoryByteLimit)

	avgDiskUtilization    := getAvgFloat64(diskUtilizationSlice)
	aboveCeilingNumber := getGreaterThanStone(diskUtilizationSlice, 0.85)

	fmt.Printf("++++++++++++++++++++++++++++++++++++\n")
	fmt.Printf("[INFO] %v\n", time.Now())

	printStatefulSetState(stsInfoGlobal)

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
			pvStatistics := pvInfos.GetStatisticsByPVName(pvName)
			fmt.Println("len(DiskIOPS): ", len(pvStatistics.DiskIOPS),
							", len(DiskWriteMBPS): ", len(pvStatistics.DiskWriteMBPS),
							", len(DiskUtilization): ", len(pvStatistics.DiskUtilization),
							", len(DiskReadMBPS): ", len(pvStatistics.DiskReadMBPS))
			fmt.Printf("%-40s \n", pvName)
			diskIOPS, err        := pvStatistics.GetLastDiskIOPS()
			if err != nil {
				log.Fatal("getLastDiskIOPS: ", err)
			}
			diskReadMBPS, err    := pvStatistics.GetLastDiskReadMBPS()
			if err != nil {
				log.Fatal("getLastDiskReadMBPS: ", err)
			}
			diskWriteMBPS, err   := pvStatistics.GetLastDiskWriteMBPS()
			if err != nil {
				log.Fatal("getLastWriteMBPS: ", err)
			}
			diskUtilization, err := pvStatistics.GetLastDiskUtilization()
			if err != nil {
				log.Fatal("getLastDiskUtilization: ", err)
			}
			fmt.Printf("From pv_collector: diskIOPS: %-10.6f, diskReadMBPS: %-10.6f, diskWriteMBPS: %-10.6f, diskUtilization: %-10.6f\n\n",
				diskIOPS, diskReadMBPS, diskWriteMBPS, diskUtilization)
		}
	}
}
