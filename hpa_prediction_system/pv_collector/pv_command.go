package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

/*
 * struct command, which represent a linux shell wrapper;
 * struct command has a method: execute, which is used to execute the strcmd;
 */
type Command struct{
	CmdPath string
}
func (c *Command) initializeCmdPath(cmdPath string) {
	c.CmdPath = cmdPath
}
func (c *Command) execute(args string) (string, error) {
	if c.CmdPath == "" {
		c.CmdPath = "/bin/bash"
	}

	cmd := exec.Command(c.CmdPath, args)
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
	args      string
}
func (p *PVCommand) getDiskUtilization() (float64, error) {
	diskUtilization, err := p.cmd.execute(p.args)
	if err != nil {
		log.Println("grepFileWithTarget warn: ", p.args, " not found!")
		return 0.0, err
	}

	diskUtilization = strings.Replace(diskUtilization, "\n", "", -1)
	if diskUtilization == "" {
		return -1.0, nil
	}
	utilization, err := strconv.ParseFloat(diskUtilization, 64)
	if err != nil {
		log.Println("strconv.Atoi error: diskUtilization: ", diskUtilization,  err)
		return 0.0, err
	}

	return utilization, err
}
func (p *PVCommand) getDiskIOPS() (float64, error) {
	diskIOPS, err := p.cmd.execute(p.args)
	if err != nil {
		log.Println("grepFileWithTarget warn: ", p.args, " not found!")
		return 0.0, err
	}

	diskIOPS = strings.Replace(diskIOPS, "\n", "", -1)
	if diskIOPS == "" {
		return -1.0, nil
	}
	iops, err := strconv.ParseFloat(diskIOPS, 64)
	if err != nil {
		log.Println("strconv.Atoi error: diskIOPS: ", diskIOPS, err)
		return 0.0, err
	}

	return iops, err
}
func (p *PVCommand) getDiskReadMBPS() (float64, error) {
	diskReadKbps, err := p.cmd.execute(p.args)
	if err != nil {
		log.Println("grepFileWithTarget warn: ", p.args, " not found!")
		return 0.0, err
	}

	diskReadKbps = strings.Replace(diskReadKbps, "\n", "", -1)
	if diskReadKbps == "" {
		return -1.0, nil
	}
	readKbps, err := strconv.ParseFloat(diskReadKbps, 64)
	if err != nil {
		log.Println("strconv.Atoi error: diskReadKbps: ", diskReadKbps, err)
		return 0.0, err
	}
	readMbps := readKbps / 1024

	return readMbps, err
}
func (p *PVCommand) getDiskWriteMBPS() (float64, error) {
	diskWriteKbps, err := p.cmd.execute(p.args)
	if err != nil {
		log.Println("grepFileWithTarget warn: ", p.args, " not found!")
		return 0.0, err
	}

	diskWriteKbps = strings.Replace(diskWriteKbps, "\n", "", -1)
	if diskWriteKbps == "" {
		return -1.0, nil
	}
	writeKbps, err := strconv.ParseFloat(diskWriteKbps, 64)
	if err != nil {
		log.Println("strconv.Atoi error: diskWriteKbps: ", diskWriteKbps, err)
		return 0.0, err
	}
	writeMbps := writeKbps / 1024

	return writeMbps, err
}
