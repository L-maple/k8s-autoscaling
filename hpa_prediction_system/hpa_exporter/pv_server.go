package main

import (
	"context"
	pb "github.com/k8s-autoscaling/hpa_prediction_system/pv_monitor"
	"strconv"
)

type server struct {
	pb.UnimplementedPVServiceServer
}

func (s *server) RequestPVNames(ctx context.Context, in *pb.PVRequest) (*pb.PVResponse, error) {
	var pvNames []string

	stsMutex.RLock()
	defer stsMutex.RUnlock()
	if stsInfoGlobal.PodInfos == nil {
		return &pb.PVResponse{PvNames: pvNames}, nil
	}

	for _, podInfo := range stsInfoGlobal.GetPodInfos() {
		for _, pvName := range podInfo.GetPVNames() {
			pvNames = append(pvNames, pvName)
		}
	}
	return &pb.PVResponse{ PvNames: pvNames }, nil
}

func (s *server) ReplyPVInfos(ctx context.Context, pvInfoRequests *pb.PVInfosRequest) (*pb.PVInfosResponse, error) {
	StrTimestamp := strconv.FormatInt(pvInfoRequests.Timestamp, 10)

	diskInfoInMemoryMutex.Lock()
	defer diskInfoInMemoryMutex.Unlock()
	for pvName, pvInfo := range pvInfoRequests.GetPVInfos() {
		StrIOPS        := strconv.FormatFloat(float64(pvInfo.PVDiskIOPS), 'f', 6, 64)
		StrUtilization := strconv.FormatFloat(float64(pvInfo.PVDiskUtilization), 'f', 6, 64)
		StrReadMBPS    := strconv.FormatFloat(float64(pvInfo.PVDiskReadMBPS), 'f', 6, 64)
		StrWriteMBPS   := strconv.FormatFloat(float64(pvInfo.PVDiskWriteMBPS), 'f', 6, 64)

		// 将PV信息添加到内存数组中
		diskIOPSInMemory[pvName] = append(diskIOPSInMemory[pvName], []string{StrTimestamp, StrIOPS})
		if len(diskIOPSInMemory[pvName]) > DiskInfoInMemoryNumber {
			startIndex := len(diskIOPSInMemory[pvName]) - DiskInfoInMemoryNumber
			diskIOPSInMemory[pvName] = diskIOPSInMemory[pvName][startIndex:]
		}

		diskUtilizationInMemory[pvName] = append(diskUtilizationInMemory[pvName], []string{StrTimestamp, StrUtilization})
		if len(diskUtilizationInMemory[pvName]) > DiskInfoInMemoryNumber {
			startIndex := len(diskUtilizationInMemory[pvName]) - DiskInfoInMemoryNumber
			diskUtilizationInMemory[pvName] = diskUtilizationInMemory[pvName][startIndex:]
		}

		diskReadMBPSInMemory[pvName] = append(diskReadMBPSInMemory[pvName], []string{StrTimestamp, StrReadMBPS})
		if len(diskReadMBPSInMemory[pvName]) > DiskInfoInMemoryNumber {
			startIndex := len(diskReadMBPSInMemory[pvName]) - DiskInfoInMemoryNumber
			diskReadMBPSInMemory[pvName] = diskReadMBPSInMemory[pvName][startIndex:]
		}

		diskWriteMBPSInMemory[pvName] = append(diskWriteMBPSInMemory[pvName], []string{StrTimestamp, StrWriteMBPS})
		if len(diskWriteMBPSInMemory[pvName]) > DiskInfoInMemoryNumber {
			startIndex := len(diskWriteMBPSInMemory[pvName]) - DiskInfoInMemoryNumber
			diskWriteMBPSInMemory[pvName] = diskWriteMBPSInMemory[pvName][startIndex:]
		}
		//fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		//fmt.Println(pvName)
		//fmt.Println(pvInfo.PVDiskIOPS, pvInfo.PVDiskUtilization, pvInfo.PVDiskReadMBPS, pvInfo.PVDiskWriteMBPS)
		//fmt.Println("===================================")
	}
	return &pb.PVInfosResponse{Status: 1}, nil
}
