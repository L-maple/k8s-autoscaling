package main

import (
	"testing"
)

func TestGetDiskWriteMBPSWithTarget(t *testing.T) {
	target := "centos-root"
	pvCmd := PVCommand{Command{}, target}

	if writeMbps, err := pvCmd.getDiskWriteMBPS(target); err != nil {
		t.Error("pvCmd.getDiskWriteMBPS error, target: ", target)
	} else {
		t.Log("DiskWriteMBPS: ", writeMbps)
	}
}

func TestGetDiskWriteMBPSWithoutTarget(t *testing.T) {
	target := "hello-world"
	pvCmd := PVCommand{Command{}, target}

	if writeMbps, err := pvCmd.getDiskWriteMBPS(target); err != nil {
		t.Error("pvCmd.getDiskWriteMBPS error, target: ", target)
	} else {
		t.Log("DiskWriteMBPS: ", writeMbps)
	}
}

func TestGetDiskReadMBPSWithTarget(t *testing.T) {
	target := "centos-root"
	pvCmd := PVCommand{Command{}, target}

	if readMbps, err := pvCmd.getDiskReadMBPS(target); err != nil {
		t.Error("pvCmd.getDiskReadMBPS error, target: ", target)
	} else {
		t.Log("DiskReadMBPS: ", readMbps)
	}
}

func TestGetDiskReadMBPSWithoutTarget(t *testing.T) {
	target := "hello_world"
	pvCmd := PVCommand{Command{}, target}

	if readMbps, err := pvCmd.getDiskReadMBPS(target); err != nil {
		t.Error("pvCmd.getDiskReadMBPS error, target: ", target)
	} else {
		t.Log("DiskReadMBPS: ", readMbps)
	}
}

func TestGetDiskIOPSWithTarget(t *testing.T) {
	target := "centos-root"
	pvCmd := PVCommand{Command{}, target}

	if iops, err := pvCmd.getDiskIOPS(target); err != nil {
		t.Error("pvCmd.getDiskIOPS error, target: ", target)
	} else {
		t.Log("DiskIOPS: ", iops)
	}
}

func TestGetDiskIOPSWithoutTarget(t *testing.T) {
	target := "hello_world"
	pvCmd := PVCommand{Command{}, target}

	if iops, err := pvCmd.getDiskIOPS(target); err != nil {
		t.Error("pvCmd.getDiskIOPS error, target: ", target)
	} else {
		t.Log("DiskIOPS: ", iops)
	}
}

func TestGetDiskUtilizationWithTarget(t *testing.T) {
	target := "centos-root"
	pvCmd := PVCommand{Command{}, target}

	if utilization, err := pvCmd.getDiskUtilization(target); err != nil {
		t.Error("pvCmd.getDiskUtilization error, target: ", target)
	} else {
		t.Log("DiskUtilization: ", utilization)
	}
}

func TestGetDiskUtilizationWithoutTarget(t *testing.T) {
	target := "hello_world"
	pvCmd := PVCommand{Command{}, target}

	if utilization, err := pvCmd.getDiskIOPS(target); err != nil {
		t.Error("pvCmd.getDiskUtilization error, target: ", target)
	} else {
		t.Log("DiskUtilization: ", utilization)
	}
}