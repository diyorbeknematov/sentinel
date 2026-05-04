package cli

import (
	"fmt"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config values",
	Run: func(cmd *cobra.Command, args []string) {
		// Mavjud configni yukla (bo'lmasa yangi yaratadi)
		cfg, err := config.Load("config.yaml")
		if err != nil {
			cfg = &config.Config{}
		}

		// Faqat berilgan flaglarni yangilaydi
		if apiKey != "" {
			cfg.APIKey = apiKey
		}
		if server != "" {
			cfg.ServerURL = server
		}

		if err := config.Save("config.yaml", cfg); err != nil {
			fmt.Println("❌ Saqlashda xatolik:", err)
			return
		}

		fmt.Println("✅ Config saqlandi")
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show config",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load("config.yaml")
		if err != nil {
			fmt.Println("❌ Config topilmadi")
			return
		}

		fmt.Println("\nJoriy config:")
		fmt.Println("─────────────────────────")
		fmt.Println("  API Key  :", maskKey(cfg.APIKey))
		fmt.Println("  Server   :", cfg.ServerURL)
		fmt.Println("  Agent ID :", cfg.AgentID)
		fmt.Println("─────────────────────────")
		fmt.Println()
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset config",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Save("config.yaml", &config.Config{}); err != nil {
			fmt.Println("❌ Reset xatolik:", err)
			return
		}
		fmt.Println("🧹 Config tozalandi")
	},
}

func init() {
	// Flaglar — set command uchun
	configSetCmd.Flags().StringVar(&apiKey, "api-key", "", "API key")
	configSetCmd.Flags().StringVar(&server, "server", "", "Server URL")

	// Sub-commandlarni bog'lash
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configResetCmd)

	rootCmd.AddCommand(configCmd)
}

// API keyni yashiradi: sk-abc... → sk-••••••••
func maskKey(key string) string {
	if key == "" {
		return "—"
	}
	if len(key) <= 5 {
		return "••••••"
	}
	return key[:5] + "••••••••"
}