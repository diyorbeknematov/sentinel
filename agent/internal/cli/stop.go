/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. PID fayl o'qish
		data, err := os.ReadFile(pidFilePath)
		if err != nil {
			fmt.Println("Agent ishlamayapti.")
			return
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil || pid <= 0 {
			fmt.Println("PID fayl noto'g'ri.")
			os.Remove(pidFilePath)
			return
		}

		// 2. Process topish
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println("Process topilmadi.")
			os.Remove(pidFilePath)
			return
		}

		// 3. Process ishlaydimi tekshirish
		if proc.Signal(syscall.Signal(0)) != nil {
			fmt.Println("Agent allaqachon to'xtatilgan.")
			os.Remove(pidFilePath)
			return
		}

		// 4. SIGTERM yuborish — to'xtatish
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			fmt.Println("Agentni to'xtatishda xatolik:", err)
			return
		}

		// 5. PID faylni o'chirish
		os.Remove(pidFilePath)

		fmt.Printf("Agent to'xtatildi (PID %d)\n", pid)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
