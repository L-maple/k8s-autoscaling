package main

import (
	"fmt"
	"github.com/sasha-s/go-deadlock"
	"math"
	"strconv"
	"time"
)


const (
	/* HPA Finite State*/
	FreeState      = 0
	StressState    = 10
	ScaleUpState   = 101
)


type HPAFiniteStateMachine struct {
	rwLock                   deadlock.RWMutex

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
	return h.finiteState
}
func (h *HPAFiniteStateMachine) GetTimerFlag() int {
	return h.timerFlag
}
func (h *HPAFiniteStateMachine) GetStabilizationWindowTime() int64 {
	return h.stabilizationWindowTime
}
/*
 * transferFromScaleUpToFreeState: 该方法将使得hpaState的状态从ScaleUp转到Free
 * 该方法只能由 StateTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromScaleUpToFreeState() {
	h.finiteState             = FreeState
	h.stabilizationWindowTime = math.MaxInt64
	h.timerFlag               = NoneTimerFlag
	fsmLog.Println("transferFromScaleUpToFreeState called: hpaFSM transfer to FreeState.")
}
/*
 * transferFromFreeToStressState: 该方法将使得hpaFSM的状态从Free转到Stress
 * 该方法由 cpuTimer, diskIOPSTimer, diskMBPSTimer 和 diskUtilizationTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromFreeToStressState(stabilizationWindowTime int64, timerFlag int) {
	h.finiteState             = StressState
	h.stabilizationWindowTime = stabilizationWindowTime
	h.timerFlag               = timerFlag
}
/*
 * resetStressState: 该方法重置hpaFSM, 变更稳定窗口时间 和 计时器标记
 * 该方法由 cpuTimer, diskIOPSTimer, diskMBPSTimer 和 diskUtilizationTimer 调用
 */
func (h *HPAFiniteStateMachine) resetStressState(stabilizationWindowTime int64, timerFlag int) {
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
	h.finiteState             = FreeState
	h.stabilizationWindowTime = math.MaxInt64
	h.timerFlag               = NoneTimerFlag
}
/*
 * transferFromStressToScaleUpState: 该方法将使得hpaState从Stress状态转移至ScaleUp状态
 * 该方法只能由 StateTimer 调用
 */
func (h *HPAFiniteStateMachine) transferFromStressToScaleUpState() {
	fsmLog.Println(time.Now(), "扩容原因: ", h.GetScaleUpReason())

	h.finiteState             = ScaleUpState
	h.stabilizationWindowTime = math.MaxInt64
	h.timerFlag               = NoneTimerFlag
}
/*
 * 返回稳定窗口
 */
func (h *HPAFiniteStateMachine) GetScaleUpReason() string {
	var reason string
	if h.timerFlag == DiskUtilizationTimerFlag {
		reason = "[Stress -> ScaleUp 扩容原因] DiskUtilization计时器达到稳定窗口时间~\n"
	} else if h.timerFlag == DiskIOPSTimerFlag {
		reason = "[Stress -> ScaleUp 扩容原因] DiskIOPS计时器达到稳定窗口时间~\n"
	} else if h.timerFlag == DiskMBPSTimerFlag {
		reason = "[Stress -> ScaleUp 扩容原因] DiskMBPS计时器达到稳定窗口时间~\n"
	} else if h.timerFlag == CPUTimerFlag {
		reason = "[Stress -> ScaleUp 扩容原因] CPU计时器达到稳定窗口时间~\n"
	} else {
		reason = "[Stress -> ScaleUp 扩容原因] 扩容原因未知: " + strconv.Itoa(h.timerFlag) + "~\n"
	}

	pvNumbers := fmt.Sprintf("~~~~~~~~~~~pvNumbers from pvInfos: %d~~~~~~~~~\n", pvInfos.GetPVNumbers())
	reason += pvNumbers

	return reason
}

func getLatestDiskMetricsInfo() string {
	iops := fmt.Sprintf("%f", pvInfos.GetAvgLastDiskIOPS())
	readMBPS := fmt.Sprintf("%f", pvInfos.GetAvgLastDiskReadMBPS())
	writeMBPS := fmt.Sprintf("%f", pvInfos.GetAvgLastDiskWriteMBPS())
	utilization := fmt.Sprintf("%f", pvInfos.GetAvgLastDiskUtilization())

	metricsInfo := "[系统指标信息] {disk_utilizaiton}: " + utilization +
		"; {disk_readMBPS}: " + readMBPS +
		"; {disk_writeMBPS}: " + writeMBPS +
		"; {disk_iops}: " + iops + "; \n"

	return metricsInfo
}

// TODO: 补充下CPU和内存的指标信息打印