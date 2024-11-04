package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func monitorProcess(process *Process) {
	for {
		printProcessInfo(process, 0)
		time.Sleep(500 * time.Millisecond)

		// Check if the process still exists
		if !processExists(process.PID) {
			break
		}
	}
}

func printProcessInfo(process *Process, space int) {
	clearConsole()
	printProcessInfoRecursive(process, space, !showOnly)
}

func printProcessInfoRecursive(process *Process, space int, printChildren bool) {
	filePath := fmt.Sprintf("/proc/%d/status", process.PID)
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	printSeparator(space)
	fmt.Printf("%*sPID: %d\n", space, "", process.PID)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			switch fields[0] {
			case "Name:":
				fmt.Printf("%*sName: %s\n", space+2, "", fields[1])
			case "State:":
				fmt.Printf("%*sState: %s\n", space+2, "", fields[1])
			case "VmRSS:":
				mem, _ := strconv.ParseFloat(fields[1], 64)
				mem /= 1024.0
				fmt.Printf("%*sVmRSS: %.2f MB\n", space+2, "", mem)
			case "Threads:":
				ppid, _ := strconv.Atoi(fields[1])
				fmt.Printf("%*sThreads: %d\n", space+2, "", ppid)
			case "PPid:":
				ppid, _ := strconv.Atoi(fields[1])
				fmt.Printf("%*sParent PID: %d\n", space+2, "", ppid)
			case "Command:":
				fmt.Printf("%*sCommand: %s\n", space+2, "", strings.Join(fields[1:], " "))
			case "Uid:":
				uid, _ := strconv.Atoi(fields[1])
				fmt.Printf("%*sUser: %d\n", space+2, "", uid)
			}
		}
	}

	if printChildren {
		children := getChildProcesses(process.PID)
		for _, childPID := range children {
			childProcess := createProcess(childPID)
			printProcessInfoRecursive(childProcess, space+2, true)
		}
	}
}

func printSeparator(space int) {
	fmt.Printf("%*s%s\n", space, "", strings.Repeat("-", 40))
}

func clearConsole() {
	cmd := exec.Command("clear") // for Linux/Mac
	cmd.Stdout = os.Stdout
	cmd.Run()
}
