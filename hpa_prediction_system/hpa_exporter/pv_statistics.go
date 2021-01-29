package main

import "strconv"

func getDiskIOPSs(pvName string) [][]string {
	diskInfoInMemoryMutex.RLock()
	defer diskInfoInMemoryMutex.RUnlock()

	return diskIOPSInMemory[pvName]
}

func getLastDiskIOPS(pvName string) (float64, error) {
	diskIOPSSlice := getDiskIOPSs(pvName)
	if len(diskIOPSSlice) == 0 {
		return 0.0, nil
	}

	StrDiskIOPS := diskIOPSSlice[len(diskIOPSSlice)-1][1]

	return strconv.ParseFloat(StrDiskIOPS, 32)
}

func getDiskWriteMBPSs(pvName string) [][]string {
	diskInfoInMemoryMutex.RLock()
	defer diskInfoInMemoryMutex.RUnlock()

	return diskWriteMBPSInMemory[pvName]
}

func getLastWriteMBPS(pvName string) (float64, error) {
	writeMBPSSlice := getDiskWriteMBPSs(pvName)
	if len(writeMBPSSlice) == 0 {
		return 0.0, nil
	}

	StrDiskWriteMBPS := writeMBPSSlice[len(writeMBPSSlice)-1][1]

	return strconv.ParseFloat(StrDiskWriteMBPS, 32)
}

func getDiskReadMBPSs(pvName string) [][]string {
	diskInfoInMemoryMutex.RLock()
	defer diskInfoInMemoryMutex.RUnlock()

	return diskReadMBPSInMemory[pvName]
}

func getLastDiskReadMBPS(pvName string) (float64, error) {
	diskReadMBPSSlice := getDiskReadMBPSs(pvName)
	if len(diskReadMBPSSlice) == 0 {
		return 0.0, nil
	}

	StrDiskReadMBPS := diskReadMBPSSlice[len(diskReadMBPSSlice)-1][1]

	return strconv.ParseFloat(StrDiskReadMBPS, 32)
}

func getDiskUtilizations(pvName string) [][]string {
	diskInfoInMemoryMutex.RLock()
	defer diskInfoInMemoryMutex.RUnlock()

	return diskUtilizationInMemory[pvName]
}

func getLastDiskUtilization(pvName string) (float64, error) {
	diskUtilizationSlice := getDiskUtilizations(pvName)
	if len(diskUtilizationSlice) == 0 {
		return 0.0, nil
	}

	StrDiskUtilization := diskUtilizationSlice[len(diskUtilizationSlice)-1][1]

	return strconv.ParseFloat(StrDiskUtilization, 32)
}