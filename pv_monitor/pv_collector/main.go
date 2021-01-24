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

	pvCmd.cmd.initializeCmdPath(diskUtilizationScript)
	diskUtilization, err := pvCmd.getDiskUtilization()
	if err != nil {
		log.Fatal("pvCmd.getDiskUtilization: ", err)
	}
	fmt.Println(diskUtilization)

	// FIXME: 对target进行处理才能调用iostat命令
	// 比如: lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf05 的 设备名为 centos-lvm--43c80d34--b593--4f7d--b7bf--a45c8f4fdf05
	pvCmd.cmd.initializeCmdPath(diskIOPSScript)
	diskIOPS, err := pvCmd.getDiskIOPS()
	if err != nil {
		log.Fatal("pvCmd.getDiskIOPS: ", err)
	}
	fmt.Println(diskIOPS)

	pvCmd.cmd.initializeCmdPath(diskReadKbpsScript)
	diskReadMbps, err := pvCmd.getDiskReadMBPS()
	if err != nil {
		log.Fatal("pvCmd.getDiskReadMBPS: ", err)
	}
	fmt.Println(diskReadMbps)

	pvCmd.cmd.initializeCmdPath(diskWriteKbpsScript)
	diskWriteMbps, err := pvCmd.getDiskWriteMBPS()
	if err != nil {
		log.Fatal("pvCmd.getDiskWriteMBPS: ", err)
	}
	fmt.Println(diskWriteMbps)
}

func init() {
	flag.IntVar(&intervalTime, "s", 10, "collector interval")
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
			//TODO: 先判定target是否存在于文件系统中
			// ...

			// 对target的指标信息进行处理
			handlePVMetricsWithScripts(target)
		}

		var pvInfos map[string]*pb.PVInfo
		fmt.Println(time.Now(), "sendPVMetrics...")
		sendPVMetrics(pvServiceClient, pvInfos)
		fmt.Println(time.Now(), ", this client send pvInfos to Server successfully~")

		time.Sleep(time.Duration(intervalTime) * time.Second)
	}
}
