package middleware

import (
	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/bot/session"

	tb "gopkg.in/telebot.v3"
)

var superAdminIDs = map[int64]bool{
	6012322301: true,
	7405650419: true,
}

func IsChatAdmin() tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			if superAdminIDs[c.Sender().ID] {
				return next(c)
			}
			if !chat.IsChatAdmin(c.Bot(), c.Chat(), c.Sender().ID) {
				return c.Reply("您不是当前会话的管理员")
			}
			v := c.Get(session.StoreKeyMentionChat.String())
			if v != nil {
				mentionChat, ok := v.(*tb.Chat)
				if !ok {
					return c.Reply("内部错误")
				}
				if !chat.IsChatAdmin(c.Bot(), mentionChat, c.Sender().ID) {
					return c.Reply("您不是当前会话的管理员")
				}
			}
			return next(c)
		}
	}
}
