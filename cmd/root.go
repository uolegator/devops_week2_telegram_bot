package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	token       string
	mode        string // "polling" or "webhook" (polling by default)
	webhookURL  string
	webhookHost string
	webhookPort int
	adminID     int64
)

var rootCmd = &cobra.Command{
	Use:   "tgbot",
	Short: "Telegram bot with Cobra + Telebot",
	Long:  "A functional Telegram bot starter with root command, settings, and message handlers.",
}

func init() {
	// Global flags (can also be set via env)
	rootCmd.PersistentFlags().StringVar(&token, "token", os.Getenv("TELEGRAM_TOKEN"), "Telegram Bot API token (or set TELEGRAM_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&mode, "mode", getEnv("BOT_MODE", "polling"), "Run mode: polling|webhook")
	rootCmd.PersistentFlags().StringVar(&webhookURL, "webhook-url", os.Getenv("WEBHOOK_URL"), "Public HTTPS URL for webhook")
	rootCmd.PersistentFlags().StringVar(&webhookHost, "webhook-host", getEnv("WEBHOOK_HOST", "0.0.0.0"), "Webhook bind host")
	rootCmd.PersistentFlags().IntVar(&webhookPort, "webhook-port", getEnvInt("WEBHOOK_PORT", 8080), "Webhook bind port")
	rootCmd.PersistentFlags().Int64Var(&adminID, "admin-id", getEnvInt64("ADMIN_ID", 0), "Admin user ID (optional)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var out int
		fmt.Sscanf(v, "%d", &out)
		if out != 0 {
			return out
		}
	}
	return def
}

func getEnvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		var out int64
		fmt.Sscanf(v, "%d", &out)
		if out != 0 {
			return out
		}
	}
	return def
}
