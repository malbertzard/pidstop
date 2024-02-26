package main

import (
	"os/exec"
	"strconv"
	"strings"
)

type Process struct {
	PID       int
	Name      string
	State     string
	VmRSS     float64
	PPID      int
	Command   string
	User      int
	ChildPIDs []int
}

func createProcess(pid int) *Process {
	return &Process{PID: pid}
}

func getChildProcesses(parentPid int) []int {
	cmd := exec.Command("pgrep", "-P", strconv.Itoa(parentPid))
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var children []int
	pids := strings.Fields(string(output))
	for _, pidStr := range pids {
		pid, err := strconv.Atoi(pidStr)
		if err == nil {
			children = append(children, pid)
		}
	}

	return children
}

func getPIDFromName(name string) int {
	cmd := exec.Command("pgrep", name)
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	pidStr := strings.TrimSpace(string(output))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0
	}

	return pid
}
