package main

import (
	"context"
	pb "github.com/k8s-autoscaling/pv_monitor/pv_monitor"
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

func (s *server) ReplyPVInfos(ctx context.Context, in *pb.PVInfosRequest) (*pb.PVInfosResponse, error) {
	return &pb.PVInfosResponse{Status: 1}, nil
}
