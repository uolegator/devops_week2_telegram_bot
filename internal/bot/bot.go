package bot

import (
	"fmt"
	"time"

	tb "gopkg.in/telebot.v4"
)

func NewBot(settings tb.Settings, adminID int64) (*tb.Bot, error) {
	// Optional: add custom client, verbose logs, etc.
	b, err := tb.NewBot(settings)
	if err != nil {
		return nil, fmt.Errorf("init bot: %w", err)
	}

	// Example: set bot commands shown in Telegram UI
	err = b.SetCommands([]tb.Command{
		{Text: "start", Description: "Start the bot"},
		{Text: "help", Description: "Show help"},
		{Text: "settings", Description: "Open settings"},
	})
	if err != nil {
		// non-fatal
		b.Logger().Println("SetCommands warning:", err)
	}

	// Example: custom error handler (optional)
	b.Use(func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			if err := next(c); err != nil {
				b.Logger().Println("handler error:", err)
			}
			return nil
		}
	})

	return b, nil
}

// Utility: common reply options
func defaultReplyOpts() *tb.SendOptions {
	return &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}
}

// Example: building a settings keyboard
func SettingsKeyboard() *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{ResizeKeyboard: true}
	btnAbout := m.Data("ℹ️ About", "about_btn")
	btnNotif := m.Data("🔔 Toggle notifications", "notif_btn")
	btnLang := m.Data("🌐 Language", "lang_btn")

	m.Inline(
		m.Row(btnAbout),
		m.Row(btnNotif),
		m.Row(btnLang),
	)
	return m
}

func Now() string {
	return time.Now().Format(time.RFC3339)
}
