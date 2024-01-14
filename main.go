package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "process-monitor",
		Short: "Monitor the VMRSS and other information of a process and its children",
		Run:   runMonitor,
	}
	processName string
	command     string
	showOnly    bool
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().IntP("pid", "p", 0, "PID of the process to monitor")
	rootCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process to monitor")
	rootCmd.Flags().StringVarP(&command, "command", "c", "", "Command to run and monitor")
	rootCmd.Flags().BoolVarP(&showOnly, "show-only", "s", false, "Show only the entry process (exclude children)")
	viper.BindPFlag("pid", rootCmd.Flags().Lookup("pid"))
	viper.BindPFlag("name", rootCmd.Flags().Lookup("name"))
	viper.BindPFlag("command", rootCmd.Flags().Lookup("command"))
	viper.BindPFlag("show-only", rootCmd.Flags().Lookup("show-only"))

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()
}

func initConfig() {
	if viper.Get("pid") == nil {
		viper.Set("pid", 0)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runMonitor(cmd *cobra.Command, args []string) {
	pid := viper.GetInt("pid")

	if pid == 0 && processName == "" && command == "" {
		fmt.Println("Please provide a valid PID using the '--pid' flag, a process name using the '--name' flag, or a command using the '--command' flag.")
		os.Exit(1)
	}

	if command != "" {
		pid = runCommand(command)
		if pid == 0 {
			fmt.Printf("Failed to run command: %s\n", command)
			os.Exit(1)
		}
	}

	if pid == 0 {
		pid = getPIDFromName(processName)
		if pid == 0 {
			fmt.Printf("Process with name '%s' not found.\n", processName)
			os.Exit(1)
		}
	}

	for {
		printProcessInfo(pid, 0)
		time.Sleep(500 * time.Millisecond)

		// Check if the process still exists
		if !processExists(pid) {
			break
		}
	}
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

func runCommand(cmdStr string) int {
	cmd := exec.Command("sh", "-c", cmdStr)
	err := cmd.Start()
	if err != nil {
		return 0
	}

	return cmd.Process.Pid
}

func printProcessInfo(pid, space int) {
	clearConsole()
	printProcessInfoRecursive(pid, space, !showOnly)
}

func printProcessInfoRecursive(pid, space int, printChildren bool) {
	filePath := fmt.Sprintf("/proc/%d/status", pid)
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	printSeparator(space)
	fmt.Printf("%*sPID: %d\n", space, "", pid)

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
		children := getChildProcesses(pid)
		for _, child := range children {
			printProcessInfoRecursive(child, space+2, true)
		}
	}
}

func printSeparator(space int) {
	fmt.Printf("%*s%s\n", space, "", strings.Repeat("-", 40))
}

func processExists(pid int) bool {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid))
	err := cmd.Run()
	return err == nil
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

func clearConsole() {
	cmd := exec.Command("clear") // for Linux/Mac
	cmd.Stdout = os.Stdout
	cmd.Run()
}
