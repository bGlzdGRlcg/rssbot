package middleware

import (
	"github.com/indes/flowerss-bot/internal/bot/session"

	tb "gopkg.in/telebot.v3"
)

func IsChatAdmin() tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			v := c.Get(session.StoreKeyMentionChat.String())
			if v != nil {
				mentionChat, ok := v.(*tb.Chat)
			}
			return next(c)
		}
	}
}
