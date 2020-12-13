package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
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

func main() {
	cmd := command{}

	if result, err := cmd.execute("lvdisplay"); err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err)
		return
	}
}
