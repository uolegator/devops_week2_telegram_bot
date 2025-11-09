package bot

import (
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v4"
)

func RegisterHandlers(b *tb.Bot, adminID int64) {
	// --- Command handlers ---
	b.Handle("/start", onStart(b))
	b.Handle("/help", onHelp(b))
	b.Handle("/settings", onSettings(b))

	// --- Content-type handlers ---
	b.Handle(tb.OnText, onText(b))
	b.Handle(tb.OnPhoto, onPhoto(b))
	b.Handle(tb.OnSticker, onSticker(b))
	b.Handle(tb.OnDocument, onDocument(b))
	b.Handle(tb.OnInlineQuery, onInlineQuery(b)) // optional, demonstrates type-specific handling

	// --- Callback handlers (for settings keyboard) ---
	b.Handle(&tb.Callback{Unique: "about_btn"}, func(c tb.Context) error {
		return c.Respond(&tb.CallbackResponse{Text: "Telegram bot starter powered by telebot.v4 + Cobra."})
	})
	b.Handle(&tb.Callback{Unique: "notif_btn"}, func(c tb.Context) error {
		return c.Edit("🔔 Notifications toggled (demo).")
	})
	b.Handle(&tb.Callback{Unique: "lang_btn"}, func(c tb.Context) error {
		return c.Edit("🌐 Language settings (demo).")
	})
}

// --- Commands ---

func onStart(b *tb.Bot) func(c tb.Context) error {
	return func(c tb.Context) error {
		user := c.Sender()
		msg := fmt.Sprintf(
			"*Welcome, %s!* 👋\n\nUse /help to see available commands.\n_Time:_ `%s`",
			displayName(user), Now(),
		)
		return c.Send(msg, defaultReplyOpts())
	}
}

func onHelp(b *tb.Bot) func(c tb.Context) error {
	return func(c tb.Context) error {
		help := "*Commands*\n" +
			"/start - Start the bot\n" +
			"/help - Show this help\n" +
			"/settings - Open settings\n\n" +
			"*Tips*\n" +
			"- Send me text, photos, stickers, or documents.\n" +
			"- I’ll respond based on message type and content."
		return c.Send(help, defaultReplyOpts())
	}
}

func onSettings(b *tb.Bot) func(c tb.Context) error {
	return func(c tb.Context) error {
		return c.Send("⚙️ Settings", &tb.SendOptions{
			ParseMode:   tb.ModeMarkdown,
			ReplyMarkup: SettingsKeyboard(),
		})
	}
}

// --- Content-type routing ---

func onText(b *tb.Bot) func(c tb.Context) error {
	return func(c tb.Context) error {
		text := strings.TrimSpace(c.Text())

		// Simple intent detection
		switch {
		case text == "":
			return nil
		case strings.EqualFold(text, "hi") || strings.EqualFold(text, "hello"):
			return c.Send("👋 Hello! Type /help for options.")
		case strings.Contains(strings.ToLower(text), "time"):
			return c.Send("⏰ Current time: " + Now())
		default:
			// Echo with minimal formatting
			reply := fmt.Sprintf("You said:\n```\n%s\n```", text)
			return c.Send(reply, defaultReplyOpts())
		}
	}
}

func onPhoto(b *tb.Bot) func(c tb.Context) error {
	return func(c tb.Context) error {
		ph := c.Message().Photo
		caption := c.Message().Caption
		msg := fmt.Sprintf("📷 Nice photo! (fileID: %s)\nCaption: %s", ph.File.FileID, caption)
		return c.Send(msg)
	}
}

func onSticker(b *tb.Bot) func(c tb.Context) error {
	return func(c tb.Context) error {
		st := c.Message().Sticker
		return c.Send(fmt.Sprintf("😄 Cool sticker! (emoji: %s, set: %s)", st.Emoji, st.SetName))
	}
}

func onDocument(b *tb.Bot) func(c tb.Context) error {
	return func(c tb.Context) error {
		doc := c.Message().Document
		return c.Send(fmt.Sprintf("📄 Got a document: %s (%d bytes)", doc.FileName, doc.FileSize))
	}
}

func onInlineQuery(b *tb.Bot) func(c tb.Context) error {
	return func(c tb.Context) error {
		// Minimal example: answer inline query with a single article.
		q := c.Query()
		if q == nil {
			return nil
		}
		result := &tb.ArticleResult{
			Title:       "Echo",
			Description: "Echo back your query",
			Text:        "You typed: " + q.Text,
		}
		result.SetResultID("echo-1")
		return c.Answer(&tb.QueryResponse{
			Results:   tb.Results{result},
			CacheTime: 1,
		})
	}
}

func displayName(u *tb.User) string {
	if u == nil {
		return "there"
	}
	if u.FirstName != "" || u.LastName != "" {
		return strings.TrimSpace(u.FirstName + " " + u.LastName)
	}
	if u.Username != "" {
		return "@" + u.Username
	}
	return "there"
}
