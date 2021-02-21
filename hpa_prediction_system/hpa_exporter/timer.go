package main

import (
	rs "github.com/k8s-autoscaling/hpa_prediction_system/hpa_exporter/resource_statistics"
	"log"
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
		if hpaFSM.GetState() == StressState &&
			hpaFSM.GetStabilizationWindowTime() >= time.Now().Unix() {
			if hpaFSM.GetTimerFlag() == NoneTimerFlag {
				log.Fatal("transferFromStressToScaleUpState Error: timerFlag is NoneTimerFlag")
			}
			hpaFSM.transferFromStressToScaleUpState()
		}

		// 状态从 ScaleUp 到 Free
		currentPodNumber := len(stsInfoGlobal.GetPodNames())
		if previousPodNumber < currentPodNumber {
			scaleUpFinished = true
		}

		// TODO: 能否判定扩容完成了？测试验证下
		// TODO: 可能存在协程死锁的问题，需要排查下
		if hpaFSM.GetState() == ScaleUpState && scaleUpFinished == true {
			fsmLog.Println("##StateTimer## transferFromScaleUpToFreeState: ",
								"hpaFSM.GetState: ", hpaFSM.finiteState,
								"hpaFSM.GetTimerFlag: ", hpaFSM.timerFlag,
								"hpaFSM.GetStabilizationWindowTime: ", hpaFSM.stabilizationWindowTime)
			hpaFSM.transferFromScaleUpToFreeState()
			scaleUpFinished = false
		}

		previousPodNumber = currentPodNumber

		time.Sleep(time.Duration(20) * time.Second)

		fsmLog.Println("##StateTimer## FSMState:", hpaFSM.finiteState,
							"stabilizationWindowTime: ", hpaFSM.stabilizationWindowTime,
							"timerFlag: ", hpaFSM.timerFlag)
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
	return (podCounter - aboveCeilingNumber < ReplicasAmount) || (avgDiskUtilization >= 0.7)
}
func (d *DiskUtilizationTimer) Run() {
	d.stabilizationWindowTime = 0

	for {
		// 状态从 Free 到 Stress
		podNameAndInfo := stsInfoGlobal.GetPodInfos()
		podCounter := len(podNameAndInfo)
		if podCounter == 0 {   // 说明stsInfoGlobal中还没有统计信息
			fsmLog.Println("##DiskUtilizationTimer## podCounter is zero...")
			time.Sleep(time.Duration(intervalTime) * time.Second)
			continue
		}

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
			stabilizationWindowTime := time.Now().Unix() + 60  // 1分钟稳定窗口时间
			if hpaFSM.GetState() == FreeState {
				fsmLog.Println("##DiskUtilizationTimer## transferFromFreeToStressState: ",
									"podCounter: ", podCounter,
									"aboveCeilingNumber: ", aboveCeilingNumber,
									"avgDiskUtilization: ", avgDiskUtilization)
				hpaFSM.transferFromFreeToStressState(stabilizationWindowTime, DiskUtilizationTimerFlag)
				d.SetStabilizationWindowTime(stabilizationWindowTime)
			}
			if hpaFSM.GetState() == StressState {
				if hpaFSM.GetStabilizationWindowTime() > stabilizationWindowTime {
					fsmLog.Println("##DiskUtilizationTimer## resetStressState: ",
						"podCounter: ", podCounter,
						"aboveCeilingNumber: ", aboveCeilingNumber,
						"avgDiskUtilization: ", avgDiskUtilization)
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
				fsmLog.Println("##DiskUtilizationTimer## transferFromStressToFreeState: ",
					"podCounter: ", podCounter,
					"aboveCeilingNumber: ", aboveCeilingNumber,
					"avgDiskUtilization: ", avgDiskUtilization)
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
