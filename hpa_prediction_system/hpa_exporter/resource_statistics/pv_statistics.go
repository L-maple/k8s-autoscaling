package statistics

import (
	"log"
	"strconv"
)

type PVInfos struct {
	NameAndStatistics map[string]PVStatistics
}

func (p *PVInfos) Initialize() {
	p.NameAndStatistics = make(map[string]PVStatistics)
}

func (p PVInfos) GetStatisticsByPVName(pvName string) PVStatistics {
	return p.NameAndStatistics[pvName]
}

func (p PVInfos) SetStatisticsByPVName(pvName string, statistics PVStatistics) {
	p.NameAndStatistics[pvName] = statistics
}

func (p PVInfos) GetAvgLastDiskIOPS() float64 {
	totalLastDiskIOPS, number := 0.0, 0
	for _, statistics := range p.NameAndStatistics {
		iops, err := statistics.GetLastDiskIOPS()
		if err != nil {
			log.Fatal("statistics.GetLastDiskIOPS: ", err)
		}
		totalLastDiskIOPS += iops
		number++
	}
	return totalLastDiskIOPS / float64(number)
}

func (p PVInfos) GetAvgLastDiskReadMBPS() float64 {
	totalLastDiskReadMBPS, number := 0.0, 0
	for _, statistics := range p.NameAndStatistics {
		iops, err := statistics.GetLastDiskReadMBPS()
		if err != nil {
			log.Fatal("statistics.GetLastDiskMBPS: ", err)
		}
		totalLastDiskReadMBPS += iops
		number++
	}
	return totalLastDiskReadMBPS / float64(number)
}

func (p PVInfos) GetAvgLastDiskWriteMBPS() float64 {
	totalLastDiskWriteMBPS, number := 0.0, 0
	for _, statistics := range p.NameAndStatistics {
		mbps, err := statistics.GetLastDiskWriteMBPS()
		if err != nil {
			log.Fatal("statistics.GetLastDiskWriteMBPS: ", err)
		}
		totalLastDiskWriteMBPS += mbps
		number++
	}
	return totalLastDiskWriteMBPS / float64(number)
}

func (p PVInfos) GetAvgLastDiskUtilization() float64 {
	totalLastDiskUtilization, number := 0.0, 0
	for _, statistics := range p.NameAndStatistics {
		utilization, err := statistics.GetLastDiskUtilization()
		if err != nil {
			log.Fatal("statistics.GetLastUtilization: ", err)
		}
		totalLastDiskUtilization += utilization
		number++
	}
	return totalLastDiskUtilization / float64(number)
}

/************************************************/

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

func (p PVStatistics)GetLastDiskWriteMBPS() (float64, error) {
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
