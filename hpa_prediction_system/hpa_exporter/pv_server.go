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

	//fmt.Println(":::::len(pvInfoRequests.GetPVInfos): ", len(pvInfoRequests.GetPVInfos()))
	for pvName, pvInfo := range pvInfoRequests.GetPVInfos() {
		StrIOPS        := strconv.FormatFloat(float64(pvInfo.PVDiskIOPS), 'f', 6, 64)
		StrUtilization := strconv.FormatFloat(float64(pvInfo.PVDiskUtilization), 'f', 6, 64)
		StrReadMBPS    := strconv.FormatFloat(float64(pvInfo.PVDiskReadMBPS), 'f', 6, 64)
		StrWriteMBPS   := strconv.FormatFloat(float64(pvInfo.PVDiskWriteMBPS), 'f', 6, 64)

		// 将PV信息添加到内存数组中
		pvStatistics := pvInfos.GetStatisticsByPVName(pvName)
		// IOPS
		pvStatistics.DiskIOPS = append(pvStatistics.DiskIOPS, []string{StrTimestamp, StrIOPS})
		if len(pvStatistics.DiskIOPS) > DiskInfoInMemoryNumber {
			startIndex := len(pvStatistics.DiskIOPS) - DiskInfoInMemoryNumber
			pvStatistics.DiskIOPS = pvStatistics.DiskIOPS[startIndex:]
		}
		// Utilization
		pvStatistics.DiskUtilization = append(pvStatistics.DiskUtilization, []string{StrTimestamp, StrUtilization})
		if len(pvStatistics.DiskUtilization) > DiskInfoInMemoryNumber {
			startIndex := len(pvStatistics.DiskUtilization) - DiskInfoInMemoryNumber
			pvStatistics.DiskUtilization = pvStatistics.DiskUtilization[startIndex:]
		}
		// ReadMbps
		pvStatistics.DiskReadMBPS = append(pvStatistics.DiskReadMBPS, []string{StrTimestamp, StrReadMBPS})
		if len(pvStatistics.DiskReadMBPS) > DiskInfoInMemoryNumber {
			startIndex := len(pvStatistics.DiskReadMBPS) - DiskInfoInMemoryNumber
			pvStatistics.DiskReadMBPS = pvStatistics.DiskReadMBPS[startIndex:]
		}
		// WriteMbps
		pvStatistics.DiskWriteMBPS = append(pvStatistics.DiskWriteMBPS, []string{StrTimestamp, StrWriteMBPS})
		if len(pvStatistics.DiskWriteMBPS) > DiskInfoInMemoryNumber {
			startIndex := len(pvStatistics.DiskWriteMBPS) - DiskInfoInMemoryNumber
			pvStatistics.DiskWriteMBPS = pvStatistics.DiskWriteMBPS[startIndex:]
		}
		//fmt.Println("--------------------------------------")
		//fmt.Println("len(DiskIOPS): ", len(pvStatistics.DiskIOPS),
		//				", len(DiskReadMBPS): ", len(pvStatistics.DiskReadMBPS),
		//				", len(DiskWriteMBPS): ", len(pvStatistics.DiskWriteMBPS),
		//				", len(DiskUtilization): ", len(pvStatistics.DiskUtilization))
		//fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

		pvInfos.SetStatisticsByPVName(pvName, pvStatistics)
	}

	return &pb.PVInfosResponse{Status: 1}, nil
}
