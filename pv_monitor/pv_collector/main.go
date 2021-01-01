package main

import (
	"flag"
	"fmt"
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
type command struct{}

func (c *command) execute(cmdstr string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", cmdstr)

	/*
	 * Create the command pipe
	 */
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("execute: cmd.StdoutPipe error: ", err)
		return "", err
	}

	/*
	 * Execute the command
	 */
	if err := cmd.Start(); err != nil {
		log.Println("execute: cmd.Start error: ", err)
		return "", err
	}

	/*
	 * Read all inputs
	 */
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

func grepFileWithTarget(target string, tmpFileName string, cmd command) (string, error) {
	// 对tmpFileName文件使用grep target命令找到对应
	utilizationAndTargetCmd:= fmt.Sprintf("grep %s %s", target, tmpFileName)

	if targetUtilization, err := cmd.execute(utilizationAndTargetCmd); err != nil {
		log.Println("cmd.execute utilizationAndTargetCmd error: ", err)
		return "", err
	} else {
		return targetUtilization, nil
	}
}

func saveDfInfo(tmpFileName string, cmd command) {
	// 读取文件系统使用量信息，保存到tmpFileName中
	targetUtilizationCmd := fmt.Sprintf("df --output=pcent,target")
	if targetUtilizations, err := cmd.execute(targetUtilizationCmd); err != nil {
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
	sleepTime int
)

func init() {
	flag.IntVar(&sleepTime, "s", 15, "collector interval")
}

func main() {
	flag.Parse()

	cmd := command{}

	for {
		tmpFileName := "targetUtilization.txt"
		target := "/var/lib/docker"

		saveDfInfo(tmpFileName, cmd)

		if utilizationAndTarget, err := grepFileWithTarget(target, tmpFileName, cmd); err != nil {		
			log.Println("grepFileWithTarget error: ", err)
			return
		} else {
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
			fmt.Println(target, " ", float64(utilization) / 100.0)
		}

		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}


