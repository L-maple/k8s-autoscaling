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
	scaleUpTime              int64
}
func (h *HPAFiniteStateMachine) Initialize() {
	h.finiteState = FreeState
	h.scaleUpTime = math.MaxInt64
	h.stabilizationWindowTime = math.MaxInt64
}
func (h *HPAFiniteStateMachine) GetState() int {
	h.stateMutex.RLock()
	defer h.stateMutex.RUnlock()

	currentState := h.finiteState
	return currentState
}
func (h *HPAFiniteStateMachine) GetStabilizationWindowTime() int64 {
	h.stateMutex.RLock()
	defer h.stateMutex.RUnlock()

	stabilizationWindowTime := h.stabilizationWindowTime
	return stabilizationWindowTime
}
func (h *HPAFiniteStateMachine) GetScaleUpTime() int64 {
	h.stateMutex.RLock()
	defer h.stateMutex.RUnlock()

	scaleUpTime := h.scaleUpTime
	return scaleUpTime
}

/*
 * transferFromScaleUpToFreeState: 该方法将使得hpaState的状态从ScaleUp转到Free
 * 该方法只能由 StateTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromScaleUpToFreeState() {
	h.stateMutex.Lock()
	h.stateMutex.Unlock()
	h.finiteState = FreeState
	h.scaleUpTime = math.MaxInt64   // TODO: 对这里的-1还需要重新思考, 如何保证状态迁移的正常进行，同时扩容时能将状态转到FreeState
	fmt.Println("transferFromScaleUpToFreeState called: hpaFSM transfer to FreeState.")
}
/*
 * transferFromFreeToStressState: 该方法将使得hpaState的状态从Free转到Stress
 * 该方法只能由 cpuTimer, diskIOPSTimer, diskMBPSTimer 和 diskUtilizationTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromFreeToStressState(stabilizationWindowTime int64, scaleUpTime int64) {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()
	h.finiteState = StressState
	h.stabilizationWindowTime = stabilizationWindowTime
	h.scaleUpTime = scaleUpTime
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
	h.scaleUpTime = math.MaxInt64
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
	h.scaleUpTime = math.MaxInt64
}

type StateTimer struct {}
func (s StateTimer) Run() {
	for {
		if hpaFSM.GetState() == StressState && hpaFSM.GetScaleUpTime() >= time.Now().Unix() {
			hpaFSM.transferFromStressToScaleUpState()
		}

		if hpaFSM.GetState() == ScaleUpState && { // TODO: 如何判定扩容完成了？
			hpaFSM.transferFromScaleUpToFreeState()
		}

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskUtilizationTimer struct {}
func (d DiskUtilizationTimer) Run() {
	for {
		stsMutex.RLock()
		podNameAndInfo := stsInfoGlobal.GetPodInfos()
		stsMutex.RUnlock()

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
		aboveCeilingNumber := getGreaterThanStone(diskUtilizationSlice, 0.85)

		if podCounter-aboveCeilingNumber < ReplicasAmount || avgDiskUtilization >= 0.85 {
			if hpaFSM.GetState() == FreeState {
				stabilizationWindowTime, scaleUpTime := time.Now().Unix(), time.Now().Unix()
				hpaFSM.transferFromFreeToStressState(stabilizationWindowTime, scaleUpTime)
			}
		}
		// TODO: 增加从Stress到Free的逻辑

		time.Sleep(time.Duration(5) * time.Second)
	}
}

type CPUTimer struct {}
func (c CPUTimer) Run() {
	for {
		// TODO: 需要完成逻辑
		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskIOPSTimer struct {}
func (d DiskIOPSTimer) Run() {
	for {
		// TODO: 需要完成逻辑
		time.Sleep(time.Duration(5) * time.Second)
	}
}

type DiskMBPSTimer struct {}
func (d DiskMBPSTimer) Run() {
	for {
		// TODO: 需要完成逻辑
		time.Sleep(time.Duration(5) * time.Second)
	}
}