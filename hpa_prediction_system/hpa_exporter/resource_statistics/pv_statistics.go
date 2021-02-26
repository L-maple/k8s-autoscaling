package statistics

import (
	"github.com/sasha-s/go-deadlock"
	"log"
	"strconv"
	"time"
)

type PVInfos struct {
	RwLock            deadlock.RWMutex
	NameAndStatistics map[string]PVStatistics
}

func (p *PVInfos) Initialize() {
	p.NameAndStatistics = make(map[string]PVStatistics)
}

func (p *PVInfos) GetPVNumbers() int {
	return len(p.NameAndStatistics)
}

func (p *PVInfos) GetStatisticsByPVName(pvName string) PVStatistics {
	return p.NameAndStatistics[pvName]
}

func (p *PVInfos) SetStatisticsByPVName(pvName string, statistics PVStatistics) {
	p.NameAndStatistics[pvName] = statistics
}

func (p *PVInfos) GetAvgLastDiskIOPS(PVNames []string) float64 {
	totalLastDiskIOPS, number := 0.0, 0
	for pvName, statistics := range p.NameAndStatistics {
		if _, res := Find(PVNames, pvName); res == false {
			continue
		}
		iops, err := statistics.GetLastDiskIOPS()
		if err != nil {
			log.Fatal("statistics.GetLastDiskIOPS: ", err)
		}
		totalLastDiskIOPS += iops
		number++
	}

	if number == 0 {
		return 0.0
	}
	return totalLastDiskIOPS / float64(number)
}
func (p *PVInfos) GetAvgLastRangeDiskIOPS(timeRange int64) float64 {
	// TODO: 添加函数
	return 0.0
}

func (p *PVInfos) GetAvgLastDiskReadMBPS(PVNames []string) float64 {
	totalLastDiskReadMBPS, number := 0.0, 0
	for pvName, statistics := range p.NameAndStatistics {
		if _, res := Find(PVNames, pvName); res == false {
			continue
		}
		iops, err := statistics.GetLastDiskReadMBPS()
		if err != nil {
			log.Fatal("statistics.GetLastDiskMBPS: ", err)
		}
		totalLastDiskReadMBPS += iops
		number++
	}
	if number == 0 {
		return 0.0
	}
	return totalLastDiskReadMBPS / float64(number)
}

func (p *PVInfos) GetAvgLastDiskWriteMBPS(PVNames []string) float64 {
	totalLastDiskWriteMBPS, number := 0.0, 0
	for pvName, statistics := range p.NameAndStatistics {
		if _, res := Find(PVNames, pvName); res == false {
			continue
		}
		mbps, err := statistics.GetLastDiskWriteMBPS()
		if err != nil {
			log.Fatal("statistics.GetLastDiskWriteMBPS: ", err)
		}
		totalLastDiskWriteMBPS += mbps
		number++
	}
	if number == 0 {
		return 0.0
	}
	return totalLastDiskWriteMBPS / float64(number)
}

func (p *PVInfos) GetAvgLastDiskUtilization(PVNames []string) float64 {
	totalLastDiskUtilization, number := 0.0, 0
	for pvName, statistics := range p.NameAndStatistics {
		if _, res := Find(PVNames, pvName); res == false {
			continue
		}
		utilization, err := statistics.GetLastDiskUtilization()
		if err != nil {
			log.Fatal("statistics.GetLastUtilization: ", err)
		}
		totalLastDiskUtilization += utilization
		number++
	}
	if number == 0 {
		return 0
	}

	return totalLastDiskUtilization / float64(number)
}

func (p *PVInfos) GetAvgLastDiskUtilizationTest(PVNames []string) (float64,[]float64) {
	totalLastDiskUtilization, number := 0.0, 0
	var utilizations []float64
	for pvName, statistics := range p.NameAndStatistics {
		if _, res := Find(PVNames, pvName); res == false {
			continue
		}
		utilization, err := statistics.GetLastDiskUtilization()
		if err != nil {
			log.Fatal("statistics.GetLastUtilization: ", err)
		}
		utilizations = append(utilizations, utilization)
		totalLastDiskUtilization += utilization
		number++
	}
	if number == 0 {
		return 0, []float64{}
	}

	return totalLastDiskUtilization / float64(number), utilizations
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

func (p PVStatistics)GetTimeDurationAvgDiskIOPS(timeDurationSec int64) float64 {
	diskIOPSSlice := p.GetDiskIOPSs()
	if len(diskIOPSSlice) == 0 {
		return 0.0
	}

	startTimeStamp := time.Now().Unix() - timeDurationSec
	totalIOPS, number := 0.0, 0
	for i := len(diskIOPSSlice)-1; i >= 0; i-- {
		StrTimeStamp, StrDiskIOPS := diskIOPSSlice[i][0], diskIOPSSlice[i][1]
		timestamp, err := strconv.ParseInt(StrTimeStamp, 10, 64)
		if err != nil {
			log.Fatal("strconv.ParseFloat: ", err)
		}
		diskIOPS, err := strconv.ParseFloat(StrDiskIOPS, 32)
		if err != nil {
			log.Fatal("strconv.ParseFloat: ", err)
		}

		if timestamp < startTimeStamp {  // 如果该记录的timestamp不在考虑范围内，那么退出循环
			break
		}

		totalIOPS += diskIOPS
		number++
	}
	avgIOPS := totalIOPS / float64(number)

	return avgIOPS
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

func (p PVStatistics)GetTimeDurationAvgDiskWriteMBPS(timeDurationSec int64) float64 {
	writeMBPSSlice := p.GetDiskWriteMBPSs()
	if len(writeMBPSSlice) == 0 {
		return 0.0
	}

	startTimeStamp := time.Now().Unix() - timeDurationSec
	totalWriteMBPS, number := 0.0, 0
	for i := len(writeMBPSSlice)-1; i >= 0; i-- {
		StrTimeStamp, StrDiskWriteMBPS := writeMBPSSlice[i][0], writeMBPSSlice[i][1]
		timestamp, err := strconv.ParseInt(StrTimeStamp, 10, 64)
		if err != nil {
			log.Fatal("strconv.ParseFloat: ", err)
		}
		writeMBPS, err := strconv.ParseFloat(StrDiskWriteMBPS, 32)
		if err != nil {
			log.Fatal("strconv.ParseFloat: ", err)
		}

		if timestamp < startTimeStamp {  // 如果该记录的timestamp不在考虑范围内，那么退出循环
			break
		}

		totalWriteMBPS += writeMBPS
		number++
	}
	avgWriteMBPS := totalWriteMBPS / float64(number)

	return avgWriteMBPS
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

func (p PVStatistics)GetTimeDurationAvgDiskReadMBPS(timeDurationSec int64) float64 {
	readMBPSSlice := p.GetDiskReadMBPSs()
	if len(readMBPSSlice) == 0 {
		return 0.0
	}

	startTimeStamp := time.Now().Unix() - timeDurationSec
	totalReadMBPS, number := 0.0, 0
	for i := len(readMBPSSlice)-1; i >= 0; i-- {
		StrTimeStamp, StrDiskReadMBPS := readMBPSSlice[i][0], readMBPSSlice[i][1]
		timestamp, err := strconv.ParseInt(StrTimeStamp, 10, 64)
		if err != nil {
			log.Fatal("strconv.ParseFloat: ", err)
		}
		readMBPS, err := strconv.ParseFloat(StrDiskReadMBPS, 32)
		if err != nil {
			log.Fatal("strconv.ParseFloat: ", err)
		}

		if timestamp < startTimeStamp {  // 如果该记录的timestamp不在考虑范围内，那么退出循环
			break
		}

		totalReadMBPS += readMBPS
		number++
	}
	avgReadMBPS := totalReadMBPS / float64(number)

	return avgReadMBPS
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

func (p PVStatistics)GetTimeDurationAvgDiskUtilization(timeDurationSec int64) float64 {
	diskUtilizationSlice := p.GetDiskUtilizations()
	if len(diskUtilizationSlice) == 0 {
		return 0.0
	}

	startTimeStamp := time.Now().Unix() - timeDurationSec
	totalDiskUtilization, number := 0.0, 0
	for i := len(diskUtilizationSlice)-1; i >= 0; i-- {
		StrTimeStamp, StrDiskUtilization := diskUtilizationSlice[i][0], diskUtilizationSlice[i][1]
		timestamp, err := strconv.ParseInt(StrTimeStamp, 10, 64)
		if err != nil {
			log.Fatal("strconv.ParseFloat: ", err)
		}
		diskUtilization, err := strconv.ParseFloat(StrDiskUtilization, 32)
		if err != nil {
			log.Fatal("strconv.ParseFloat: ", err)
		}

		if timestamp < startTimeStamp {  // 如果该记录的timestamp不在考虑范围内，那么退出循环
			break
		}

		totalDiskUtilization += diskUtilization
		number++
	}
	avgDiskUtilization := totalDiskUtilization / float64(number)

	return avgDiskUtilization
}