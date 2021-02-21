package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)


const (
	/* HPA Finite State*/
	FreeState      = 0
	StressState    = 1
	ScaleUpState   = 2
)


type HPAFiniteStateMachine struct {
	stateMutex               sync.RWMutex

	finiteState              int
	stabilizationWindowTime  int64
	timerFlag                int
}
func (h *HPAFiniteStateMachine) Initialize() {
	h.finiteState             = FreeState
	h.stabilizationWindowTime = math.MaxInt64
	h.timerFlag               = NoneTimerFlag
}
func (h *HPAFiniteStateMachine) GetState() int {
	//h.stateMutex.RLock()
	//defer h.stateMutex.RUnlock()

	return h.finiteState
}
func (h *HPAFiniteStateMachine) GetTimerFlag() int {
	return h.timerFlag
}
func (h *HPAFiniteStateMachine) GetStabilizationWindowTime() int64 {
	h.stateMutex.RLock()
	defer h.stateMutex.RUnlock()

	return h.stabilizationWindowTime
}
/*
 * transferFromScaleUpToFreeState: 该方法将使得hpaState的状态从ScaleUp转到Free
 * 该方法只能由 StateTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromScaleUpToFreeState() {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()

	h.finiteState             = FreeState
	h.stabilizationWindowTime = math.MaxInt64
	h.timerFlag               = NoneTimerFlag
	fmt.Println("transferFromScaleUpToFreeState called: hpaFSM transfer to FreeState.")
}
/*
 * transferFromFreeToStressState: 该方法将使得hpaFSM的状态从Free转到Stress
 * 该方法由 cpuTimer, diskIOPSTimer, diskMBPSTimer 和 diskUtilizationTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromFreeToStressState(stabilizationWindowTime int64, timerFlag int) {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()

	h.finiteState             = StressState
	h.stabilizationWindowTime = stabilizationWindowTime
	h.timerFlag               = timerFlag
}
/*
 * resetStressState: 该方法重置hpaFSM, 变更稳定窗口时间 和 计时器标记
 * 该方法由 cpuTimer, diskIOPSTimer, diskMBPSTimer 和 diskUtilizationTimer 调用
 */
func (h *HPAFiniteStateMachine) resetStressState(stabilizationWindowTime int64, timerFlag int) {
	h.stateMutex.Lock()
	defer h.stateMutex.Lock()

	h.finiteState             = StressState
	h.stabilizationWindowTime = stabilizationWindowTime
	h.timerFlag               = timerFlag
}
/*
 * transferFromStressToFreeState: 该方法将hpaState的状态从Stress转到Free
 * 该方法由 cpuTimer, diskIOPSTimer, diskMBPSTimer 和 diskUtilizationTimer 调用
 * TODO: 考虑下有2个timer都将状态调整到了stress，那么如何对状态进行正常操作
 */
func (h *HPAFiniteStateMachine) transferFromStressToFreeState() {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()

	h.finiteState             = FreeState
	h.stabilizationWindowTime = math.MaxInt64
	h.timerFlag               = NoneTimerFlag
}
/*
 * transferFromStressToScaleUpState: 该方法将使得hpaState从Stress状态转移至ScaleUp状态
 * 该方法只能由 StateTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromStressToScaleUpState() {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()

	h.finiteState             = ScaleUpState
	h.stabilizationWindowTime = math.MaxInt64
	h.timerFlag               = NoneTimerFlag

	log.Println(time.Now(), "扩容原因: ", h.GetScaleUpReason())
}
/*
 * 返回稳定窗口
 */
func (h *HPAFiniteStateMachine) GetScaleUpReason() string {
	h.stateMutex.RLock()
	defer h.stateMutex.RUnlock()

	if h.timerFlag == DiskUtilizationTimerFlag {
		return "[扩容] DiskUtilization计时器达到稳定窗口时间~"
	} else if h.timerFlag == DiskIOPSTimerFlag {
		return "[扩容] DiskIOPS计时器达到稳定窗口时间~"
	} else if h.timerFlag == DiskMBPSTimerFlag {
		return "[扩容] DiskMBPS计时器达到稳定窗口时间"
	} else if h.timerFlag == CPUTimerFlag {
		return "[扩容] CPU计时器达到稳定窗口时间"
	}

	return "[扩容] 扩容原因未知: " + strconv.Itoa(h.timerFlag)
}
