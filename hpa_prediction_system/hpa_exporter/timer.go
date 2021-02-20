package main

import (
	"fmt"
	rs "github.com/k8s-autoscaling/hpa_prediction_system/hpa_exporter/resource_statistics"
	"math"
	"sync"
	"time"
)

type HPAFiniteStateMachine struct {
	stateMutex               sync.RWMutex
	finiteState              int
	stabilizationWindowTime  int64
}
func (h *HPAFiniteStateMachine) Initialize() {
	h.finiteState = FreeState
	h.stabilizationWindowTime = math.MaxInt64
}
func (h *HPAFiniteStateMachine) GetState() int {
	h.stateMutex.RLock()
	defer h.stateMutex.RUnlock()

	return h.finiteState
}
func (h *HPAFiniteStateMachine) GetStabilizationWindowTime() int64 {
	h.stateMutex.RLock()
	defer h.stateMutex.RUnlock()

	stabilizationWindowTime := h.stabilizationWindowTime
	return stabilizationWindowTime
}
/*
 * transferFromScaleUpToFreeState: 该方法将使得hpaState的状态从ScaleUp转到Free
 * 该方法只能由 StateTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromScaleUpToFreeState() {
	h.stateMutex.Lock()
	h.stateMutex.Unlock()
	h.finiteState             = FreeState
	h.stabilizationWindowTime = math.MaxInt64
	fmt.Println("transferFromScaleUpToFreeState called: hpaFSM transfer to FreeState.")
}
/*
 * transferFromFreeToStressState: 该方法将使得hpaState的状态从Free转到Stress
 * 该方法只能由 cpuTimer, diskIOPSTimer, diskMBPSTimer 和 diskUtilizationTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromFreeToStressState(stabilizationWindowTime int64) {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()
	h.finiteState = StressState
	h.stabilizationWindowTime = stabilizationWindowTime
}
/*
 * TODO: 该方法该由谁来调用呢?
 * TODO: 考虑下有2个timer都将状态调整到了stress，那么如何对状态进行正常操作
 */
func (h *HPAFiniteStateMachine) transferFromStressToFreeState() {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()
	h.finiteState = FreeState
	h.stabilizationWindowTime = math.MaxInt64
}
/*
 * transferFromStressToScaleUpState: 该方法将使得hpaState从Stress状态转移至ScaleUp状态
 * 该方法只能由 StateTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromStressToScaleUpState() {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()
	h.finiteState = ScaleUpState
	h.stabilizationWindowTime = math.MaxInt64
}

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

		if hpaFSM.GetState() == ScaleUpState && scaleUpFinished == true { // TODO: 能否判定扩容完成了？测试验证下
			hpaFSM.transferFromScaleUpToFreeState()
		}

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskUtilizationTimer struct {}
func (d DiskUtilizationTimer) Run() {
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

		avgDiskUtilization := getAvgFloat64(diskUtilizationSlice)
		aboveCeilingNumber := getGreaterThanStone(diskUtilizationSlice, 0.7)
		// TODO: 增加时间序列预测的支持
		if podCounter - aboveCeilingNumber < ReplicasAmount || avgDiskUtilization >= 0.5 {
			if hpaFSM.GetState() == FreeState {
				stabilizationWindowTime := time.Now().Unix() + 60
				hpaFSM.transferFromFreeToStressState(stabilizationWindowTime)
			}
		}

		// TODO: 增加从Stress到Free的逻辑

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type CPUTimer struct {}
func (c CPUTimer) Run() {
	for {
		// TODO: 增加从Free到Stress的逻辑


		// TODO: 增加从Stress到Free的逻辑


		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskIOPSTimer struct {}
func (d DiskIOPSTimer) Run() {
	for {
		// TODO: 增加从Free到Stress的逻辑

		// TODO: 增加从Stress到Free的逻辑

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskMBPSTimer struct {}
func (d DiskMBPSTimer) Run() {
	for {
		// TODO: 增加从Free到Stress的逻辑

		// TODO: 增加从Stress到Free的逻辑

		time.Sleep(time.Duration(5) * time.Second)
	}
}