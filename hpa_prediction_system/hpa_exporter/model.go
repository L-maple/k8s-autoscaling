package main

import (
	"fmt"
	rs "github.com/k8s-autoscaling/hpa_prediction_system/hpa_exporter/resource_statistics"
	"log"
	"time"
)


func getHpaActivityState() int {
	// 如果 stsInfoGlobal还没初始化，那么直接返回 FreeState
	fmt.Println("stsInfoFlobal.initialize: ", stsInfoGlobal.isInitialized())

	if stsInfoGlobal.isInitialized() == false {
		stsInfoGlobal.rwLock.RLock()
		printStatefulSetState(stsInfoGlobal)
		stsInfoGlobal.rwLock.RUnlock()

		return hpaFSM.GetState()
	}

	printCurrentState()

	return hpaFSM.GetState()
}

func printCurrentState() {
	podNameAndInfo  := stsInfoGlobal.GetPodInfos()
	memoryByteLimit := stsInfoGlobal.GetMemoryByteLimit()

	var cpuUtilizationSlice    []float64
	var memoryUsageSlice       []int64
	for podName, _ := range podNameAndInfo {
		podStatisticsObj := rs.PodStatistics{
			Endpoint:  prometheusUrl,
			PodName:   podName,
			Namespace: namespaceName,
		}

		cpuUtilizationSlice  = append(cpuUtilizationSlice, podStatisticsObj.GetLastCpuUtilization())
		memoryUsageSlice     = append(memoryUsageSlice, podStatisticsObj.GetLastMemoryUsage())
	}

	// 得到CPU的平均使用率
	avgCpuUtilization    := getAvgFloat64(cpuUtilizationSlice)

	// 得到Memory的平均使用率
	avgMemoryUsage := getAvgInt64(memoryUsageSlice)
	avgMemoryUtilization := float64(avgMemoryUsage) / float64(memoryByteLimit)

	avgDiskUtilization    := pvInfos.GetAvgLastDiskUtilization()

	fmt.Printf("++++++++++++++++++++++++++++++++++++\n")
	fmt.Printf("[INFO] %v\n", time.Now())

	printStatefulSetState(stsInfoGlobal)

	fmt.Printf("avgCpuUtilization: %-30.6f, avgMemoryUtilization: %-30.6f, avgDiskUtilization: %-30.6f\n",
					avgCpuUtilization, avgMemoryUtilization, avgDiskUtilization)
	fmt.Printf("====================================\n\n")
}

func printStatefulSetState(stsInfo *StatefulSetInfo) {
	fmt.Printf("%-40s %-40s %-40s\n", "PodName", "PvcName", "PvName")
	stsInfo.rwLock.RLock()
	podInfos := stsInfo.GetPodInfos()
	stsInfo.rwLock.RUnlock()

	for podName, podInfo := range podInfos {
		fmt.Printf("%-40s ", podName)

		for _, pvcName := range podInfo.PVCNames {
			fmt.Printf("%-40s ", pvcName)
		}

		var diskUtilizationSlice []float64
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
			diskUtilizationSlice = append(diskUtilizationSlice, diskUtilization)
			fmt.Printf("From pv_collector: diskIOPS: %-10.6f, diskReadMBPS: %-10.6f, diskWriteMBPS: %-10.6f, diskUtilization: %-10.6f\n\n",
				diskIOPS, diskReadMBPS, diskWriteMBPS, diskUtilization)
		}
		aboveCeilingNumber := getAboveBoundaryNumber(diskUtilizationSlice, 0.85)
		fmt.Printf("pod Numbers: %d, aboveCeilingNumber: %d\n", len(podInfos), aboveCeilingNumber)
	}
}
