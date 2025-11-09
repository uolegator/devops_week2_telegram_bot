package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	tb "gopkg.in/telebot.v4"

	"github.com/spf13/cobra"
	appbot "github.com/yourname/telebot-cobra-starter/internal/bot"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the Telegram bot",
	RunE: func(cmd *cobra.Command, args []string) error {
		if token == "" {
			return fmt.Errorf("missing token: supply --token or TELEGRAM_TOKEN env")
		}

		switch mode {
		case "polling":
			return runPolling()
		case "webhook":
			return runWebhook()
		default:
			return fmt.Errorf("unknown mode %q (use polling|webhook)", mode)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runPolling() error {
	settings := tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := appbot.NewBot(settings, adminID)
	if err != nil {
		return err
	}
	appbot.RegisterHandlers(b, adminID)

	log.Println("Bot started in polling mode.")
	return b.Start()
}

func runWebhook() error {
	if webhookURL == "" {
		return fmt.Errorf("webhook mode requires --webhook-url (public HTTPS URL)")
	}

	// Telebot webhook settings
	settings := tb.Settings{
		Token: token,
		Poller: &tb.Webhook{
			Listen: fmt.Sprintf("%s:%d", webhookHost, webhookPort),
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: webhookURL,
			},
		},
	}

	b, err := appbot.NewBot(settings, adminID)
	if err != nil {
		return err
	}
	appbot.RegisterHandlers(b, adminID)

	// Optional graceful shutdown
	go func() {
		if err := b.Start(); err != nil {
			log.Printf("bot stopped: %v", err)
		}
	}()

	log.Printf("Bot webhook listening on %s:%d (public: %s)", webhookHost, webhookPort, webhookURL)

	// If you need to keep main alive in webhook mode:
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", webhookHost, webhookPort),
		ReadHeaderTimeout: 5 * time.Second,
	}
	// The telebot webhook poller internally registers handlers on net/http DefaultServeMux.
	// Start server:
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return srv.Shutdown(context.Background())
}
