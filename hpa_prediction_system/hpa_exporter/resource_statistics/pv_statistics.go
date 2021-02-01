package statistics

import (
	"strconv"
)

type PVStatistics struct {
	DiskIOPS, DiskUtilization, DiskReadMBPS, DiskWriteMBPS [][]string
}

func (p PVStatistics)GetDiskIOPSs() [][]string {

	return p.DiskIOPS
}

func (p PVStatistics)GetLastDiskIOPS() (float64, error) {
	diskIOPSSlice := p.GetDiskIOPSs()
	if len(diskIOPSSlice) == 0 {
		return 0.0, nil
	}

	StrDiskIOPS := diskIOPSSlice[len(diskIOPSSlice)-1][1]

	return strconv.ParseFloat(StrDiskIOPS, 32)
}

func (p PVStatistics)GetDiskWriteMBPSs() [][]string {
	return p.DiskWriteMBPS
}

func (p PVStatistics)GetLastWriteMBPS() (float64, error) {
	writeMBPSSlice := p.GetDiskWriteMBPSs()
	if len(writeMBPSSlice) == 0 {
		return 0.0, nil
	}

	StrDiskWriteMBPS := writeMBPSSlice[len(writeMBPSSlice)-1][1]

	return strconv.ParseFloat(StrDiskWriteMBPS, 32)
}

func (p PVStatistics)GetDiskReadMBPSs() [][]string {
	return p.DiskReadMBPS
}

func (p PVStatistics)GetLastDiskReadMBPS() (float64, error) {
	diskReadMBPSSlice := p.GetDiskReadMBPSs()
	if len(diskReadMBPSSlice) == 0 {
		return 0.0, nil
	}

	StrDiskReadMBPS := diskReadMBPSSlice[len(diskReadMBPSSlice)-1][1]

	return strconv.ParseFloat(StrDiskReadMBPS, 32)
}


func (p PVStatistics)GetDiskUtilizations() [][]string {
	return p.DiskUtilization
}

func (p PVStatistics)GetLastDiskUtilization() (float64, error) {
	diskUtilizationSlice := p.GetDiskUtilizations()
	if len(diskUtilizationSlice) == 0 {
		return 0.0, nil
	}

	StrDiskUtilization := diskUtilizationSlice[len(diskUtilizationSlice)-1][1]

	return strconv.ParseFloat(StrDiskUtilization, 32)
}
