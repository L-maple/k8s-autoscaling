package main

import (
	rs "github.com/k8s-autoscaling/hpa_prediction_system/hpa_exporter/resource_statistics"
	"log"
	"math"
	"time"
)

const (
	NoneTimerFlag            = 0
	StateTimerFlag           = 1
	DiskUtilizationTimerFlag = 2
	CPUTimerFlag             = 3
	DiskIOPSTimerFlag        = 4
	DiskMBPSTimerFlag        = 5

	TimerSleep               = 5
)

type StateTimer struct {}
func (s StateTimer) Run() {
	var previousPodNumber int
	for {
		previousPodNumber = len(stsInfoGlobal.GetPodNames())
		if previousPodNumber <= 0 {
			fsmLog.Println("##StateTimer## podNumber is zero...")
			time.Sleep(time.Duration(intervalTime) * time.Second)
		} else {
			break
		}
	}

	scaleUpFinished := false
	for {
		// 状态从 Stress 到 ScaleUp
		hpaFSM.rwLock.Lock()
		if hpaFSM.GetState() == StressState &&
			hpaFSM.GetStabilizationWindowTime() <= time.Now().Unix() {
			if hpaFSM.GetTimerFlag() == NoneTimerFlag {
				log.Fatal("transferFromStressToScaleUpState Error: timerFlag is NoneTimerFlag")
			}
			hpaFSM.transferFromStressToScaleUpState()
		}

		// 状态从 ScaleUp 到 Free
		_, PVUtilizations := pvInfos.GetAvgLastDiskUtilizationTest(stsInfoGlobal.GetPVs())
		currentPodNumber := len(PVUtilizations)
		if previousPodNumber != currentPodNumber {
			scaleUpFinished = true
		}

		if (hpaFSM.GetState() == ScaleUpState) && (scaleUpFinished == true) {
			fsmLog.Println("##StateTimer## transferFromScaleUpToFreeState: ",
								"hpaFSM.GetState: ", hpaFSM.GetState(),
								"hpaFSM.GetTimerFlag: ", hpaFSM.timerFlag,
								"hpaFSM.GetStabilizationWindowTime: ", hpaFSM.GetStabilizationWindowTime())
			hpaFSM.transferFromScaleUpToFreeState()
			scaleUpFinished = false
		}
		hpaFSM.rwLock.Unlock()

		fsmLog.Println("##StateTimer## FSMState:", hpaFSM.GetState(),
			"stabilizationWindowTime: ", hpaFSM.GetStabilizationWindowTime(),
			"timerFlag: ", hpaFSM.GetTimerFlag(),
			"previousPodNumber: ", previousPodNumber,
			"currentPodNumber: ", currentPodNumber)
		fsmLog.Println(getLatestDiskMetricsInfo())

		previousPodNumber = currentPodNumber

		time.Sleep(time.Duration(TimerSleep) * time.Second)
	}
}

type DiskUtilizationTimer struct {
	avgDiskUtilizationBoundary float64
	availabilityBottomBoundary float64
	availabilityUpperBoundary  float64

	adjustTime                 int64    /* avgDiskUtilizationBoundary调整时间 */
}
func (d *DiskUtilizationTimer) IsStress(podCounter int, aboveCeilingNumber int, avgDiskUtilization float64, state int) (bool, int64) {
	if state == FreeState {
		if podCounter - aboveCeilingNumber < ReplicasAmount {
			d.avgDiskUtilizationBoundary -= 1 / 2.0 * math.Abs(d.avgDiskUtilizationBoundary-d.availabilityBottomBoundary)
			return true, time.Now().Unix()
		} else if avgDiskUtilization >= d.avgDiskUtilizationBoundary {
			d.avgDiskUtilizationBoundary += 1 / 2.0 * math.Abs(d.availabilityUpperBoundary-d.avgDiskUtilizationBoundary)
			return true, time.Now().Unix() + 30
		}
	} else if state == StressState {
		if podCounter - aboveCeilingNumber < ReplicasAmount {
			return true, time.Now().Unix()
		}

		return true, time.Now().Unix() + 30
	}
	return false, -1
}
func (d *DiskUtilizationTimer) Run() {
	d.avgDiskUtilizationBoundary = 0.6
	d.availabilityBottomBoundary = 0.6
	d.availabilityUpperBoundary = 0.85

	for {
		// 状态从 Free 到 Stress
		stsInfoGlobal.rwLock.Lock()
		podNameAndInfo := stsInfoGlobal.GetPodInfos()
		stsInfoGlobal.rwLock.Unlock()

		podCounter := len(podNameAndInfo)
		if podCounter == 0 {   // 说明stsInfoGlobal中还没有统计信息
			fsmLog.Println("##DiskUtilizationTimer## podCounter is zero...")
			time.Sleep(time.Duration(intervalTime) * time.Second)
			continue
		}

		var diskUtilizationSlice []float64
		// TODO: 等将内存数据保存到influxdb后，换掉这里从pvInfos获取数据
		avgDiskUtilization, diskUtilizationSlice := pvInfos.GetAvgLastDiskUtilizationTest(stsInfoGlobal.GetPVs())
		aboveCeilingNumber := getAboveBoundaryNumber(diskUtilizationSlice, d.availabilityUpperBoundary)
		// TODO: 增加时间序列预测的支持
		isStress, stabilizationWindowTime := d.IsStress(podCounter, aboveCeilingNumber, avgDiskUtilization, hpaFSM.GetState())
		if isStress {
			hpaFSM.rwLock.Lock()
			if hpaFSM.GetState() == FreeState {
				fsmLog.Println("##DiskUtilizationTimer## transferFromFreeToStressState: ",
									"podCounter: ", podCounter,
									"aboveCeilingNumber: ", aboveCeilingNumber,
									"avgDiskUtilization: ", avgDiskUtilization)
				hpaFSM.transferFromFreeToStressState(stabilizationWindowTime, DiskUtilizationTimerFlag)
			}
			if hpaFSM.GetState() == StressState {
				if hpaFSM.GetStabilizationWindowTime() > stabilizationWindowTime {
					fsmLog.Println("##DiskUtilizationTimer## resetStressState: ",
						"podCounter: ", podCounter,
						"aboveCeilingNumber: ", aboveCeilingNumber,
						"avgDiskUtilization: ", avgDiskUtilization)
					hpaFSM.resetStressState(stabilizationWindowTime, DiskUtilizationTimerFlag)
				}
			}
			hpaFSM.rwLock.Unlock()
		}

		// 从Stress到Free的逻辑
		hpaFSM.rwLock.Lock()
		isStress, _ = d.IsStress(podCounter, aboveCeilingNumber, avgDiskUtilization, hpaFSM.GetState())
		if (hpaFSM.GetState() == StressState) &&
			(hpaFSM.GetTimerFlag() == DiskUtilizationTimerFlag) &&
			 isStress == false {
				fsmLog.Println("##DiskUtilizationTimer## transferFromStressToFreeState: ",
					"podCounter: ", podCounter,
					"aboveCeilingNumber: ", aboveCeilingNumber,
					"avgDiskUtilization: ", avgDiskUtilization)
				hpaFSM.transferFromStressToFreeState()
			}
		hpaFSM.rwLock.Unlock()

		fsmLog.Println("Dynamic Boundary: availabilityBottomBoundary: ", d.availabilityBottomBoundary,
						"availabilityUpperBoundary: ", d.availabilityUpperBoundary,
			"avgDiskUtilizationBoundary: ", d.avgDiskUtilizationBoundary)
		time.Sleep(time.Duration(TimerSleep) * time.Second)
	}
}

type CPUTimer struct {}
func (c *CPUTimer) GetStressCondition(avgCpuUtilizationFor10Min,
	avgCpuUtilizationFor20Min,
	avgCpuUtilizationFor30Min,
	diskUtilization float64) bool {

	if avgCpuUtilizationFor10Min >= 0.8 && diskUtilization >= 0.70 ||
		avgCpuUtilizationFor20Min >= 0.8 && diskUtilization >= 0.65 ||
		avgCpuUtilizationFor30Min >= 0.8 && diskUtilization >= 0.60 {
	    	return true
	}

	return false
}

func (c *CPUTimer) Run() {
	for {
		// 状态从 Free 到 Stress
		stsInfoGlobal.rwLock.Lock()
		podNameAndInfo := stsInfoGlobal.GetPodInfos()
		stsInfoGlobal.rwLock.Unlock()

		podCounter := len(podNameAndInfo)
		if podCounter == 0 {   // 说明stsInfoGlobal中还没有统计信息
			fsmLog.Println("##DiskUtilizationTimer## podCounter is zero...")
			time.Sleep(time.Duration(intervalTime) * time.Second)
			continue
		}

		var cpuUtilizationSliceFor10Min, cpuUtilizationSliceFor20Min, cpuUtilizationSliceFor30Min []float64
		for podName, _ := range podNameAndInfo {
			podStatisticsObj := rs.PodStatistics{
				Endpoint:  prometheusUrl,
				PodName:   podName,
				Namespace: namespaceName,
			}

			cpuUtilizationSliceFor10Min = append(cpuUtilizationSliceFor10Min, podStatisticsObj.GetAvgLastRangeCPUUtilization(10 * 60))
			cpuUtilizationSliceFor20Min = append(cpuUtilizationSliceFor20Min, podStatisticsObj.GetAvgLastRangeCPUUtilization(20 * 60))
			cpuUtilizationSliceFor30Min = append(cpuUtilizationSliceFor30Min, podStatisticsObj.GetAvgLastRangeCPUUtilization(30 * 60))
		}

		avgCpuUtilizationFor10Min := getAvgFloat64(cpuUtilizationSliceFor10Min)
		avgCpuUtilizationFor20Min := getAvgFloat64(cpuUtilizationSliceFor20Min)
		avgCpuUtilizationFor30Min := getAvgFloat64(cpuUtilizationSliceFor30Min)
		diskUtilization := pvInfos.GetAvgLastDiskUtilization(stsInfoGlobal.GetPVs())

		// TODO: 增加时间序列预测的支持
		if c.GetStressCondition(avgCpuUtilizationFor10Min,
								avgCpuUtilizationFor20Min,
								avgCpuUtilizationFor30Min,
								diskUtilization) == true {
			stabilizationWindowTime := time.Now().Unix() + 60  // 1分钟稳定窗口时间

			hpaFSM.rwLock.Lock()
			if hpaFSM.GetState() == FreeState {
				fsmLog.Println("##CPUTimer## transferFromFreeToStressState: ",
					"podCounter: ", podCounter,
					"avgCpuUtilizationFor10Min: ", avgCpuUtilizationFor10Min,
					"avgCpuUtilizationFor20Min: ", avgCpuUtilizationFor20Min,
					"avgCpuUtilizationFor30Min: ", avgCpuUtilizationFor30Min,
					"avgDiskUtilization: ", diskUtilization)
				hpaFSM.transferFromFreeToStressState(stabilizationWindowTime, CPUTimerFlag)
			}
			if hpaFSM.GetState() == StressState {
				if hpaFSM.GetStabilizationWindowTime() > stabilizationWindowTime {
					fsmLog.Println("##CPUTimer## resetStressState: ",
						"podCounter: ", podCounter,
						"avgCpuUtilizationFor10Min: ", avgCpuUtilizationFor10Min,
						"avgCpuUtilizationFor20Min: ", avgCpuUtilizationFor20Min,
						"avgCpuUtilizationFor30Min: ", avgCpuUtilizationFor30Min,
						"avgDiskUtilization: ", diskUtilization)
					hpaFSM.resetStressState(stabilizationWindowTime, CPUTimerFlag)
				}
			}
			hpaFSM.rwLock.Unlock()
		}

		// 从 Stress 到 Free 的逻辑
		hpaFSM.rwLock.Lock()
		if (hpaFSM.GetState() == StressState) &&
			(hpaFSM.GetTimerFlag() == DiskUtilizationTimerFlag) &&
			(c.GetStressCondition(avgCpuUtilizationFor10Min,
								  avgCpuUtilizationFor20Min,
								  avgCpuUtilizationFor30Min,
								  diskUtilization) == false) {
			fsmLog.Println("##CPUTimer## transferFromStressToFreeState: ",
				"podCounter: ", podCounter,
				"avgCpuUtilizationFor10Min: ", avgCpuUtilizationFor10Min,
				"avgCpuUtilizationFor20Min: ", avgCpuUtilizationFor20Min,
				"avgCpuUtilizationFor30Min: ", avgCpuUtilizationFor30Min,
				"avgDiskUtilization: ", diskUtilization)
			hpaFSM.transferFromStressToFreeState()
		}
		hpaFSM.rwLock.Unlock()

		time.Sleep(time.Duration(TimerSleep) * time.Second)
	}
}

type DiskIOPSTimer struct {}
func (d *DiskIOPSTimer)GetStressCondition() bool {
	// TODO: 对disk_iops设置stress条件
	return true
}
func (d *DiskIOPSTimer) Run() {
	for {
		// 状态从 Free 到 Stress
		stsInfoGlobal.rwLock.Lock()
		podNameAndInfo := stsInfoGlobal.GetPodInfos()
		stsInfoGlobal.rwLock.Unlock()

		podCounter := len(podNameAndInfo)
		if podCounter == 0 {   // 说明stsInfoGlobal中还没有统计信息
			fsmLog.Println("##DiskUtilizationTimer## podCounter is zero...")
			time.Sleep(time.Duration(intervalTime) * time.Second)
			continue
		}

		// TODO:完善 pvInfos.GetAvgLastRangeDiskIOPS(int64)
		avgDiskIOPSFor10Min := pvInfos.GetAvgLastRangeDiskIOPS(10 * 60)
		avgDiskIOPSFor20Min := pvInfos.GetAvgLastRangeDiskIOPS(20 * 60)
		avgDiskIOPSFor30Min := pvInfos.GetAvgLastRangeDiskIOPS(30 * 60)
		diskUtilization := pvInfos.GetAvgLastDiskUtilization(stsInfoGlobal.GetPVs())

		// TODO: 增加时间序列预测的支持
		if d.GetStressCondition() == true {
			stabilizationWindowTime := time.Now().Unix() + 60  // 1分钟稳定窗口时间

			hpaFSM.rwLock.Lock()
			if hpaFSM.GetState() == FreeState {
				fsmLog.Println("##CPUTimer## transferFromFreeToStressState: ",
					"podCounter: ", podCounter,
					"avgDiskIOPSFor10Min: ", avgDiskIOPSFor10Min,
					"avgDiskIOPSFor20Min: ", avgDiskIOPSFor20Min,
					"avgDiskIOPSFor30Min: ", avgDiskIOPSFor30Min,
					"avgDiskUtilization: ", diskUtilization)
				hpaFSM.transferFromFreeToStressState(stabilizationWindowTime, CPUTimerFlag)
			}
			if hpaFSM.GetState() == StressState {
				if hpaFSM.GetStabilizationWindowTime() > stabilizationWindowTime {
					fsmLog.Println("##CPUTimer## resetStressState: ",
						"podCounter: ", podCounter,
						"avgDiskIOPSFor10Min: ", avgDiskIOPSFor10Min,
						"avgDiskIOPSFor20Min: ", avgDiskIOPSFor20Min,
						"avgDiskIOPSFor30Min: ", avgDiskIOPSFor30Min,
						"avgDiskUtilization: ", diskUtilization)
					hpaFSM.resetStressState(stabilizationWindowTime, CPUTimerFlag)
				}
			}
			hpaFSM.rwLock.Unlock()
		}

		// 从 Stress 到 Free 的逻辑
		hpaFSM.rwLock.Lock()
		if (hpaFSM.GetState() == StressState) &&
			(hpaFSM.GetTimerFlag() == DiskUtilizationTimerFlag) &&
			(d.GetStressCondition() == false) {
			fsmLog.Println("##CPUTimer## transferFromStressToFreeState: ",
				"podCounter: ", podCounter,
				"avgDiskIOPSFor10Min: ", avgDiskIOPSFor10Min,
				"avgDiskIOPSFor20Min: ", avgDiskIOPSFor20Min,
				"avgDiskIOPSFor30Min: ", avgDiskIOPSFor30Min,
				"avgDiskUtilization: ", diskUtilization)
			hpaFSM.transferFromStressToFreeState()
		}
		hpaFSM.rwLock.Unlock()

		time.Sleep(time.Duration(TimerSleep) * time.Second)
	}
}

type DiskMBPSTimer struct {}
func (d *DiskMBPSTimer) Run() {
	for {
		// TODO: 增加从Free到Stress的逻辑

		// TODO: 增加从Stress到Free的逻辑

		time.Sleep(time.Duration(TimerSleep) * time.Second)
	}
}
