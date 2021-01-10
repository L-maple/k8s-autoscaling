package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/k8s-autoscaling/pv_monitor/pv_monitor"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	/* interval time */
	intervalTime int
	timeout      int

	/* server address */
	serverAddress  = "localhost:30002"

	/* script location */
	diskUtilizationScript = "./scripts/disk_utilization.sh"
	diskIOPSScript        = "./scripts/disk_iops.sh"
	diskReadKbpsScript    = "./scripts/disk_read_kbps.sh"
	diskWriteKbpsScript   = "./scripts/disk_write_kbps.sh"
)

func getPVServiceClient() (pb.PVServiceClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewPVServiceClient(conn)

	return client, conn
}

func getTargetsFromServer(pvServiceClient pb.PVServiceClient) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(intervalTime) * time.Second)
	defer cancel()
	resp, err := pvServiceClient.RequestPVNames(ctx, &pb.PVRequest{Id: "1"})
	if err != nil {
		log.Println("pvServiceClient.RequestPVNames error: ", err)
		time.Sleep(time.Duration(intervalTime) * time.Second)
		return []string{}, err
	}
	targets := resp.PvNames
	fmt.Println("targets from grpc: ", targets)

	return targets, nil
}

func sendPVMetrics(pvServiceClient pb.PVServiceClient, pvInfos map[string]*pb.PVInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(intervalTime) * time.Second)
	defer cancel()

	resp, err := pvServiceClient.ReplyPVInfos(ctx, &pb.PVInfosRequest{
		PVInfos: pvInfos,
	})
	if err != nil {
		log.Println("pvServiceClient.PVInfosRequest error: ", err)
		return
	}

	log.Println("resp.Status is ", resp.Status)
}

func handlePVMetricsWithScripts(target string) {
	pvCmd := PVCommand{Command{}, target}

	diskUtilization, err := pvCmd.getDiskUtilization(diskUtilizationScript)
	if err != nil {
		log.Fatal("pvCmd.getDiskUtilization: ", err)
	}
	fmt.Println(diskUtilization)

	diskIOPS, err := pvCmd.getDiskIOPS(diskIOPSScript)
	if err != nil {
		log.Fatal("pvCmd.getDiskIOPS: ", err)
	}
	fmt.Println(diskIOPS)

	diskReadKbps, err := pvCmd.getDiskReadKBPS(diskReadKbpsScript)
	if err != nil {
		log.Fatal("pvCmd.getDiskReadKBPS: ", err)
	}
	fmt.Println(diskReadKbps)

	diskWriteKbps, err := pvCmd.getDiskWriteKBPS(diskWriteKbpsScript)
	if err != nil {
		log.Fatal("pvCmd.getDiskWriteKBPS: ", err)
	}
	fmt.Println(diskWriteKbps)
}

func init() {
	flag.IntVar(&intervalTime, "s", 15, "collector interval")
	flag.IntVar(&timeout, "timeout", 5, "rpc request timeout")
}

func main() {
	flag.Parse()

	pvServiceClient, requestConn:= getPVServiceClient()
	defer requestConn.Close()

	for {
		targets, err := getTargetsFromServer(pvServiceClient)
		if err != nil {
			log.Fatal("getTargetsFromGrpc error: ", err)
		}

		for _, target := range targets {
			handlePVMetricsWithScripts(target)
		}

		var pvInfos map[string]*pb.PVInfo
		fmt.Println(time.Now(), "sendPVMetrics...")
		sendPVMetrics(pvServiceClient, pvInfos)
		fmt.Println(time.Now(), ", this client send pvInfos to Server successfully~")

		time.Sleep(time.Duration(intervalTime) * time.Second)
	}
}
