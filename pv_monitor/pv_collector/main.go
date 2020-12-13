package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
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
		log.Println("execute: cmd.StdoutPipe error")
		return "", err
	}

	/*
	 * Execute the command
	 */
	if err := cmd.Start(); err != nil {
		log.Println("execute: cmd.Start error")
		return "", err
	}

	/*
	 * Read all inputs
	 */
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("execute: ioutil.ReadAll error")
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		log.Println("execute: execute: cmd.Wait error")
		return "", err
	}

	return string(bytes), nil
}

type pvparser struct{}

func (pp *pvparser) GetVgName(lvstr string) (string, error) {
	if lvstr == "" {
		return "", errors.New("GetVgName: lvstr is invalid")
	}
	lvinfo := strings.Split(lvstr, ":")
	fmt.Println("vginfo: ", lvinfo)
	return lvinfo[1], nil
}

func (pp *pvparser) GetPvPath(lvstr string) (string, error) {
	if lvstr == "" {
		return "", errors.New("GetPvPath: lvstr is invalid")
	}
	lvinfo := strings.Split(lvstr, ":")

	return lvinfo[0], nil
}

func (pp *pvparser) GetPvName(lvstr string) (string, error) {
	path, err := pp.GetPvPath(lvstr)
	if err != nil {
		return "", err
	}
	paths := strings.Split(path, "/")
	return paths[len(paths)-1], nil
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
	parser := pvparser{}

	for {
		if result, err := cmd.execute("lvdisplay -c"); err == nil {
			lvs := strings.Split(result, "\n")
			for i := 0; i < len(lvs); i++ {
				fmt.Println(lvs[i])
				// infos := strings.Split(lvs[i], ":")
				pvname, err := parser.GetPvName(lvs[i])
				if err != nil {
					fmt.Println(err)
				}

				vgname, err := parser.GetVgName(lvs[i])
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(pvname, vgname)
			}
		}

		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
