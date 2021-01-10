package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/k8s-autoscaling/pv_monitor/pv_monitor"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*
 * struct command, which represent a linux shell wrapper;
 * struct command has a method: execute, which is used to execute the strcmd;
 */
type Command struct{}

func (c *Command) execute(cmdstr string, target string) (string, error) {
	cmd := exec.Command("/bin/bash", cmdstr, target)

	/* Create the command pipe */
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("execute: cmd.StdoutPipe error: ", err)
		return "", err
	}

	/* Execute the command */
	if err := cmd.Start(); err != nil {
		log.Println("execute: cmd.Start error: ", err)
		return "", err
	}

	/* Read all inputs */
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("execute: ioutil.ReadAll error: ", err)
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		log.Println("execute: execute: cmd.Wait error: ", err)
		return "", err
	}

	return string(bytes), nil
}

func grepFileWithTarget(target string, tmpFileName string, cmd Command) (string, error) {
	// 对tmpFileName文件使用grep target命令找到对应
	utilizationAndTargetCmd:= fmt.Sprintf("grep %s %s", target, tmpFileName)

	if targetUtilization, err := cmd.execute(utilizationAndTargetCmd, ""); err != nil {
		log.Println("cmd.execute utilizationAndTargetCmd warn: ", target, " not found!")
		return "", err
	} else {
		return targetUtilization, nil
	}
}

func saveDfInfo(tmpFileName string, cmd Command) {
	// 读取文件系统使用量信息，保存到tmpFileName中
	targetUtilizationCmd := fmt.Sprintf("df --output=pcent,target")
	if targetUtilizations, err := cmd.execute(targetUtilizationCmd, ""); err != nil {
		log.Println("cmd.execute targetUtilizationCmd error: ", err)
		return
	} else {
		file, err := os.Create(tmpFileName)
		if err != nil{
			log.Fatal("error: os.Create error")
		}
		defer file.Close()

		if _, err := file.WriteString(targetUtilizations); err != nil {
			log.Fatal(err.Error())
		}
	}
}

var (
	/* interval time */
	intervalTime int
	timeout      int

	/* server address */
	serverAddress  = "localhost:30002"

	/* the tmp file for pv utilization*/
	dfInfoFileName     = "df.txt"
	iostatInfoFileName = "iostat.txt"
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

func handlePVMetrics(target string) {
	cmd := Command{}
	saveDfInfo(dfInfoFileName, cmd)
	utilizationAndTarget, err := grepFileWithTarget(target, dfInfoFileName, cmd)
	if err != nil {
		log.Println("grepFileWithTarget warn: ", target, " not found!")
		return
	}
	utilizationAndTarget = strings.Trim(utilizationAndTarget, " ")
	slices := strings.Split(utilizationAndTarget, "%")
	if len(slices) <= 1 {
		log.Println("strings.Split error: ", slices)
		return
	}

	utilization, err := strconv.Atoi(slices[0])
	if err != nil {
		log.Println("strconv.Atoi error: ", err)
		return
	}
	fmt.Println(target, " ", float64(utilization)/100.0)

	// TODO: fetch all metrics to server
}

func getTargetsFromGrpc(pvServiceClient pb.PVServiceClient) ([]string, error) {
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
	cmd := Command{}

	diskUtilization, err := cmd.execute("./disk_utilization.sh", target)
	if err != nil {
		log.Println("grepFileWithTarget warn: ", target, " not found!")
		return
	}
	slices := strings.Split(diskUtilization, "\n")
	if len(slices) <= 1 {
		log.Println("strings.Split error: ", slices)
		return
	}

	utilization, err := strconv.ParseFloat(slices[0], 32)
	if err != nil {
		log.Println("strconv.Atoi error: ", err)
		return
	}
	fmt.Println(target, ": ", utilization)

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
		//targets, err := getTargetsFromGrpc(pvServiceClient)
		//if err != nil {
		//	log.Fatal("getTargetsFromGrpc error: ", err)
		//}

		targets := []string{"/sys/fs/cgroup", "/run/user/0"}
		for _, target := range targets {
			// handlePVMetrics(target)
			handlePVMetricsWithScripts(target)
		}

		var pvInfos map[string]*pb.PVInfo
		fmt.Println(time.Now(), "sendPVMetrics...")
		sendPVMetrics(pvServiceClient, pvInfos)
		fmt.Println(time.Now(), ", this client send pvInfos to Server successfully~")

		time.Sleep(time.Duration(intervalTime) * time.Second)
	}
}
