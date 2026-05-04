/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/spf13/cobra"
)

const pidFilePath = "/tmp/sentinel.pid"

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check agent status",
	
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nAgent Status:")
		fmt.Println("─────────────────────────")

		// 1. PID fayl o'qish
		data, err := os.ReadFile(pidFilePath)
		if err != nil {
			fmt.Println("  Process   : 🔴 not running")
			fmt.Println()
			return
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil || pid <= 0 {
			fmt.Println("  Process   : 🔴 not running (invalid pid)")
			fmt.Println()
			return
		}

		// 2. Process ishlaydimi tekshirish
		proc, err := os.FindProcess(pid)
		if err != nil || proc.Signal(syscall.Signal(0)) != nil {
			fmt.Println("  Process   : 🔴 not running")
			fmt.Println()
			return
		}

		fmt.Printf("  Process   : 🟢 running (PID %d)\n", pid)

		// 3. Config o'qish
		cfg, err := config.Load("config.yaml")
		if err != nil {
			fmt.Println("  Agent ID  : —")
			fmt.Println("  Server    : —")
		} else {
			fmt.Printf("  Agent ID  : %s\n", cfg.AgentID)
			fmt.Printf("  Server    : %s\n", cfg.ServerURL)
		}

		fmt.Println("─────────────────────────")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
