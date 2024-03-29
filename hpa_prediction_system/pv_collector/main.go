package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/k8s-autoscaling/hpa_prediction_system/pv_monitor"
	"github.com/sercand/kuberesolver"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
)

var (
	/* interval time */
	intervalTime int
	timeout      int

	/* server address */
	serverAddress  string

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(intervalTime * 3) * time.Second)
	defer cancel()
	resp, err := pvServiceClient.RequestPVNames(ctx, &pb.PVRequest{Id: "1"})
	if err != nil {
		log.Println("pvServiceClient.RequestPVNames error: ", err)
		return []string{}, err
	}
	targets := resp.PvNames

	return targets, nil
}

func sendPVMetrics(pvServiceClient pb.PVServiceClient, pvInfos map[string]*pb.PVInfo, timestamp int64) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(intervalTime * 3) * time.Second)
	defer cancel()

	resp, err := pvServiceClient.ReplyPVInfos(ctx, &pb.PVInfosRequest{
		PVInfos: pvInfos,
		Timestamp: timestamp,
	})
	if err != nil {
		log.Println("pvServiceClient.PVInfosRequest error: ", err)
		return -1, err
	}

	return resp.Status, nil
}

func handlePVMetricsWithScripts(target string) *pb.PVInfo {
	target = preprocess(target)
	pvCmd := PVCommand{Command{}, target}

	// get disk utilization
	pvCmd.cmd.initializeCmdPath(diskUtilizationScript)
	diskUtilization, _ := pvCmd.getDiskUtilization()

	// get disk iops
	pvCmd.cmd.initializeCmdPath(diskIOPSScript)
	diskIOPS, _ := pvCmd.getDiskIOPS()

	// get disk read mbps
	pvCmd.cmd.initializeCmdPath(diskReadKbpsScript)
	diskReadMbps, _ := pvCmd.getDiskReadMBPS()

	// get disk write mbps
	pvCmd.cmd.initializeCmdPath(diskWriteKbpsScript)
	diskWriteMbps, _ := pvCmd.getDiskWriteMBPS()

	fmt.Println("handled-target: ", target)
	fmt.Println("metrics: ", diskUtilization, diskIOPS, diskReadMbps, diskWriteMbps)

	pvInfo := pb.PVInfo{
		PVDiskUtilization: float32(diskUtilization),
		PVDiskIOPS:        float32(diskIOPS),
		PVDiskReadMBPS:    float32(diskReadMbps),
		PVDiskWriteMBPS:   float32(diskWriteMbps),
	}

	return &pvInfo
}

// 对target进行预处理
// 比如: lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf05 只保留最后的a45c8f4fdf05
// 这样对 iostat 和 df 命令都适用
func preprocess(target string) string {
	if target == "" {
		return ""
	}
	separators := strings.Split(target, "-")

	return separators[len(separators)-1]
}

func printCurrentPvInfos(targets []string, pvInfos map[string]*pb.PVInfo, status int32) {
	fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++\n")

	fmt.Printf("[INFO] %v\n", time.Now())
	fmt.Printf("Received targets are: \n")
	for index, target := range targets {
		fmt.Printf("%-60s", target)
		if (index + 1) % 2 == 0 {
			fmt.Printf("\n")
		}
	}

	fmt.Printf("\n-----------------------\n")

	if len(pvInfos) == 0 {
		fmt.Printf("pvInfos is empty.\n")
	}
	for pvName, pvInfo := range pvInfos {
		fmt.Printf("PVName: %s\n", pvName)
		fmt.Printf("Utilization: %-15.6f IOPS: %-15.6f ReadMBPS: %-15.6f WriteMBPS: %-15.6f\n",
			pvInfo.PVDiskUtilization, pvInfo.PVDiskIOPS, pvInfo.PVDiskReadMBPS, pvInfo.PVDiskWriteMBPS)
	}
	fmt.Printf("\nsend PvInfos successfully, status: %d\n", status)

	fmt.Printf("===========================================\n\n")
}

func init() {
	flag.IntVar(&intervalTime, "s", 5, "collector interval")
	flag.IntVar(&timeout, "timeout", 5, "rpc request timeout")
	//flag.StringVar(&serverAddress, "serverAddress", "localhost:30002", "hpa-exporter comm address")
	// https://pkg.go.dev/github.com/sercand/kuberesolver@v2.4.0+incompatible
	flag.StringVar(&serverAddress, "serverAddress", "kubernetes://hdfs-hpa-exporter-service.monitoring:30002/", "hpa-exporter comm address")
	//flag.StringVar(&serverAddress, "serverAddress", "kubernetes://tidb-hpa-exporter-service.monitoring:30002/", "hpa-exporter comm address")
}

func main() {
	flag.Parse()

	// https://github.com/sercand/kuberesolver/issues/16
	// register the kuberesolver builder to grpc with kubernetes schema
	kuberesolver.RegisterInCluster()

	pvServiceClient, requestConn:= getPVServiceClient()
	defer requestConn.Close()

	timestamp := time.Now().Unix()
	for {
		targets, err := getTargetsFromServer(pvServiceClient)
		if err != nil {
			log.Println("getTargetsFromGrpc error: ", err)
			time.Sleep(time.Duration(intervalTime) * time.Second)
			continue
		}

		pvInfos := make(map[string]*pb.PVInfo)
		for _, target := range targets {
			// 对target的指标信息进行处理
			pvInfo := handlePVMetricsWithScripts(target)
			// use lvdisplay to map device number(dm-*) to lvm-name
			if pvInfo.PVDiskReadMBPS < 0 || pvInfo.PVDiskWriteMBPS < 0 ||
					pvInfo.PVDiskIOPS < 0 || pvInfo.PVDiskUtilization < 0 {
				continue
			}
			pvInfos[target] = pvInfo
		}

		status, err := sendPVMetrics(pvServiceClient, pvInfos, timestamp)
		if err != nil {
			log.Println("error: ", err)
			time.Sleep(time.Duration(intervalTime) * time.Second)
			continue
		}

		printCurrentPvInfos(targets, pvInfos, status)

		nextTimestamp := timestamp + int64(intervalTime)
		for time.Now().Unix() < nextTimestamp {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
		timestamp = nextTimestamp
	}
}


