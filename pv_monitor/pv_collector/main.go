package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/k8s-autoscaling/pv_monitor/pv_monitor"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
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

type PVCommand struct{
	cmd       Command
	target    string
}
func (p *PVCommand) getDiskUtilization() (float64, error) {
	diskUtilization, err := p.cmd.execute("./disk_utilization.sh", p.target)
	if err != nil {
		log.Println("grepFileWithTarget warn: ", p.target, " not found!")
		return 0.0, err
	}
	slices := strings.Split(diskUtilization, "\n")
	if len(slices) <= 1 {
		log.Println("strings.Split error: ", slices)
		return 0.0, err
	}

	utilization, err := strconv.ParseFloat(slices[0], 32)
	if err != nil {
		log.Println("strconv.Atoi error: ", err)
		return 0.0, err
	}

	return utilization, err
}
func (p *PVCommand) getDiskIOPS() (float64, error) {
	diskIOPS, err := p.cmd.execute("./disk_iops.sh", p.target)
	if err != nil {
		log.Println("grepFileWithTarget warn: ", p.target, " not found!")
		return 0.0, err
	}
	slices := strings.Split(diskIOPS, "\n")
	if len(slices) <= 1 {
		log.Println("strings.Split error: ", slices)
		return 0.0, err
	}

	iops, err := strconv.ParseFloat(slices[0], 32)
	if err != nil {
		log.Println("strconv.Atoi error: ", err)
		return 0.0, err
	}

	return iops, err
}
func (p *PVCommand) getDiskReadKBPS()  {

}
func (p *PVCommand) getDiskWriteKBPS() {

}


var (
	/* interval time */
	intervalTime int
	timeout      int

	/* server address */
	serverAddress  = "localhost:30002"
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

	diskUtilization, err := pvCmd.getDiskUtilization()
	if err != nil {
		log.Fatal("pvCmd.getDiskUtilization: ", err)
	}
	fmt.Println(diskUtilization)

	diskIOPS, err := pvCmd.getDiskIOPS()
	if err != nil {
		log.Fatal("pvCmd.getDiskIOPS: ", err)
	}
	fmt.Println(diskIOPS)


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
