package main

import (
	"testing"
)

func TestGetDiskWriteMBPSWithTarget(t *testing.T) {
	target := "lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf05"
	pvCmd := PVCommand{Command{CmdPath: "./scripts/disk_write_kbps.sh"}, target}

	if writeMbps, err := pvCmd.getDiskWriteMBPS(); err != nil {
		t.Error("pvCmd.getDiskWriteMBPS error, target: ", target)
	} else {
		t.Log("DiskWriteMBPS: ", writeMbps)
	}
}

func TestGetDiskWriteMBPSWithoutTarget(t *testing.T) {
	target := "lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf09"
	pvCmd := PVCommand{Command{CmdPath: "./scripts/disk_write_kbps.sh"}, target}

	if writeMbps, err := pvCmd.getDiskWriteMBPS(); err != nil {
		t.Error("pvCmd.getDiskWriteMBPS error, target: ", target)
	} else {
		t.Log("DiskWriteMBPS: ", writeMbps)
	}
}

func TestGetDiskReadMBPSWithTarget(t *testing.T) {
	target := "lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf05"
	pvCmd := PVCommand{Command{CmdPath: "./scripts/disk_read_kbps.sh"}, target}

	if readMbps, err := pvCmd.getDiskReadMBPS(); err != nil {
		t.Error("pvCmd.getDiskReadMBPS error, target: ", target)
	} else {
		t.Log("DiskReadMBPS: ", readMbps)
	}
}

func TestGetDiskReadMBPSWithoutTarget(t *testing.T) {
	target := "lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf09"
	pvCmd := PVCommand{Command{CmdPath: "./scripts/disk_read_kbps.sh"}, target}

	if readMbps, err := pvCmd.getDiskReadMBPS(); err != nil {
		t.Error("pvCmd.getDiskReadMBPS error, target: ", target)
	} else {
		t.Log("DiskReadMBPS: ", readMbps)
	}
}

func TestGetDiskIOPSWithTarget(t *testing.T) {
	target := "lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf05"
	pvCmd := PVCommand{Command{CmdPath: "./scripts/disk_iops.sh"}, target}

	if iops, err := pvCmd.getDiskIOPS(); err != nil {
		t.Error("pvCmd.getDiskIOPS error, target: ", target)
	} else {
		t.Log("DiskIOPS: ", iops)
	}
}

func TestGetDiskIOPSWithoutTarget(t *testing.T) {
	target := "lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf09"
	pvCmd := PVCommand{Command{CmdPath: "./scripts/disk_iops.sh"}, target}

	if iops, err := pvCmd.getDiskIOPS(); err != nil {
		t.Error("pvCmd.getDiskIOPS error, target: ", target)
	} else {
		t.Log("DiskIOPS: ", iops)
	}
}

func TestGetDiskUtilizationWithTarget(t *testing.T) {
	target := "lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf05"
	pvCmd := PVCommand{Command{CmdPath: "./scripts/disk_utilization.sh"}, target}

	if utilization, err := pvCmd.getDiskUtilization(); err != nil {
		t.Error("pvCmd.getDiskUtilization error, target: ", target)
	} else {
		t.Log("DiskUtilization: ", utilization)
	}
}

func TestGetDiskUtilizationWithoutTarget(t *testing.T) {
	target := "lvm-43c80d34-b593-4f7d-b7bf-a45c8f4fdf09"
	pvCmd := PVCommand{Command{CmdPath: "./scripts/disk_utilization.sh"}, target}

	if utilization, err := pvCmd.getDiskIOPS(); err != nil {
		t.Error("pvCmd.getDiskUtilization error, target: ", target)
	} else {
		t.Log("DiskUtilization: ", utilization)
	}
}

