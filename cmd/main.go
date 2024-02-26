package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

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

	if processName != "" {
		pid = getPIDFromName(processName)
		if pid == 0 {
			fmt.Printf("Process with name '%s' not found.\n", processName)
			os.Exit(1)
		}
	}

	if pid == 0 {
		os.Exit(1)
		fmt.Printf("Process with pid not found.\n")
	}

	process := createProcess(pid)
	monitorProcess(process)
}

func runCommand(cmdStr string) int {
	cmd := exec.Command("sh", "-c", cmdStr)
	err := cmd.Start()
	if err != nil {
		return 0
	}

	return cmd.Process.Pid
}

func processExists(pid int) bool {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid))
	err := cmd.Run()
	return err == nil
}
