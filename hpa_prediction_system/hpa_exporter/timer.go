package main

import (
	rs "github.com/k8s-autoscaling/hpa_prediction_system/hpa_exporter/resource_statistics"
	"time"
)

const (
	NoneTimerFlag            = 0
	StateTimerFlag           = 1
	DiskUtilizationTimerFlag = 2
	CPUTimerFlag             = 3
	DiskIOPSTimerFlag        = 4
	DiskMBPSTimerFlag        = 5
)

var (
	/* 副本数量 */
	ReplicasAmount = 3
)

type StateTimer struct {}
func (s StateTimer) Run() {
	podNumber := len(stsInfoGlobal.GetPodNames())

	for {
		scaleUpFinished := false

		// 状态从 Stress 到 ScaleUp
		if hpaFSM.GetState() == StressState && hpaFSM.GetStabilizationWindowTime() >= time.Now().Unix() {
			hpaFSM.transferFromStressToScaleUpState()
		}

		// 状态从 ScaleUp 到 Free
		currentPodNumber := len(stsInfoGlobal.GetPodNames())
		if podNumber < currentPodNumber {
			scaleUpFinished = true
		}
		podNumber = currentPodNumber

		// TODO: 能否判定扩容完成了？测试验证下
		if hpaFSM.GetState() == ScaleUpState && scaleUpFinished == true {
			hpaFSM.transferFromScaleUpToFreeState()
		}

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskUtilizationTimer struct {
	stabilizationWindowTime  int64
}
func (d *DiskUtilizationTimer) SetStabilizationWindowTime(time int64) {
	d.stabilizationWindowTime = time
}
func (d *DiskUtilizationTimer) GetStabilizationWindowTime() int64 {
	return d.stabilizationWindowTime
}
func (d *DiskUtilizationTimer) GetStressCondition(podCounter int, aboveCeilingNumber int, avgDiskUtilization float64) bool {
	return podCounter - aboveCeilingNumber < ReplicasAmount || avgDiskUtilization >= 0.7
}
func (d *DiskUtilizationTimer) Run() {
	d.stabilizationWindowTime = 0

	for {
		// 状态从 Free 到 Stress
		podNameAndInfo := stsInfoGlobal.GetPodInfos()
		podCounter := len(podNameAndInfo)

		var diskUtilizationSlice []float64
		for podName, _ := range podNameAndInfo {
			podStatisticsObj := rs.PodStatistics{
				Endpoint:  prometheusUrl,
				PodName:   podName,
				Namespace: namespaceName,
			}

			diskUtilizationSlice = append(diskUtilizationSlice, podStatisticsObj.GetLastDiskUtilization())
		}
		avgDiskUtilization := pvInfos.GetAvgLastDiskUtilization()

		aboveCeilingNumber := getGreaterThanStone(diskUtilizationSlice, 0.7)
		// TODO: 增加时间序列预测的支持
		if d.GetStressCondition(podCounter, aboveCeilingNumber, avgDiskUtilization) == true {
			if hpaFSM.GetState() == FreeState {
				stabilizationWindowTime := time.Now().Unix() + 60  // 进入1分钟稳定窗口时间
				hpaFSM.transferFromFreeToStressState(stabilizationWindowTime, DiskUtilizationTimerFlag)
				d.SetStabilizationWindowTime(stabilizationWindowTime)
			}
			if hpaFSM.GetState() == StressState {
				stabilizationWindowTime := time.Now().Unix() + 60  // 进入1分钟稳定窗口时间
				if hpaFSM.GetStabilizationWindowTime() > stabilizationWindowTime {
					hpaFSM.resetStressState(stabilizationWindowTime, DiskUtilizationTimerFlag)
					d.SetStabilizationWindowTime(stabilizationWindowTime)
				}
			}
		}

		// 从Stress到Free的逻辑
		if hpaFSM.GetState() == StressState &&
			hpaFSM.GetTimerFlag() == DiskUtilizationTimerFlag &&
			hpaFSM.GetStabilizationWindowTime() < d.GetStabilizationWindowTime() {
			if d.GetStressCondition(podCounter, aboveCeilingNumber, avgDiskUtilization) == false {
				hpaFSM.transferFromStressToFreeState()
			}
		}

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type CPUTimer struct {
	stabilizationWindowTime  int64
}
func (c *CPUTimer) GetStabilizationWindowTime() int64 {
	return c.stabilizationWindowTime
}
func (c *CPUTimer) SetStabilizationWindowTime(time int64) {
	c.stabilizationWindowTime = time
}
func (c *CPUTimer) GetStressCondition() bool {
	// TODO: 添加CPU计时器时间
	return true
}
func (c *CPUTimer) Run() {
	c.stabilizationWindowTime = 0

	for {
		// TODO: 增加从Free到Stress的逻辑


		// TODO: 增加从Stress到Free的逻辑

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskIOPSTimer struct {
	stabilizationWindowTime  int
}
func (d *DiskIOPSTimer) Run() {
	d.stabilizationWindowTime = 0

	for {
		// TODO: 增加从Free到Stress的逻辑

		// TODO: 增加从Stress到Free的逻辑

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskMBPSTimer struct {
	stabilizationWindowTime  int
}
func (d *DiskMBPSTimer) Run() {
	d.stabilizationWindowTime = 0

	for {
		// TODO: 增加从Free到Stress的逻辑

		// TODO: 增加从Stress到Free的逻辑

		time.Sleep(time.Duration(5) * time.Second)
	}
}
