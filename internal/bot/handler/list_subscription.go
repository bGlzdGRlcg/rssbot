package handler

import (
	"context"
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v3"

	"github.com/bGlzdGRlcg/rssbot/internal/bot/chat"
	"github.com/bGlzdGRlcg/rssbot/internal/bot/message"
	"github.com/bGlzdGRlcg/rssbot/internal/core"
	"github.com/bGlzdGRlcg/rssbot/internal/log"
	"github.com/bGlzdGRlcg/rssbot/internal/model"
)

const (
	MaxSubsSizePerPage = 50
)

type ListSubscription struct {
	core *core.Core
}

func NewListSubscription(core *core.Core) *ListSubscription {
	return &ListSubscription{core: core}
}

func (l *ListSubscription) Command() string {
	return "/list"
}

func (l *ListSubscription) Description() string {
	return "已订阅的RSS源"
}

func (l *ListSubscription) listChatSubscription(ctx tb.Context) error {

	stdCtx := context.Background()
	sources, err := l.core.GetUserSubscribedSources(stdCtx, ctx.Chat().ID)
	if err != nil {
		log.Errorf("GetUserSubscribedSources failed, %v", err)
		return ctx.Send("获取订阅错误")
	}

	return l.replaySubscribedSources(ctx, sources)
}

func (l *ListSubscription) listChannelSubscription(ctx tb.Context, channelName string) error {
	channelChat, err := ctx.Bot().ChatByUsername(channelName)
	if err != nil {
		return ctx.Send("获取频道信息错误")
	}

	if !chat.IsChatAdmin(ctx.Bot(), channelChat, ctx.Sender().ID) {
		return ctx.Send("非频道管理员无法执行此操作")
	}

	stdCtx := context.Background()
	sources, err := l.core.GetUserSubscribedSources(stdCtx, channelChat.ID)
	if err != nil {
		log.Errorf("GetUserSubscribedSources failed, %v", err)
		return ctx.Send("获取订阅错误")
	}
	return l.replaySubscribedSources(ctx, sources)
}

func (l *ListSubscription) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	if mention != "" {
		return l.listChannelSubscription(ctx, mention)
	}
	return l.listChatSubscription(ctx)
}

func (l *ListSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}

func (l *ListSubscription) replaySubscribedSources(ctx tb.Context, sources []*model.Source) error {
	if len(sources) == 0 {
		return ctx.Send("订阅列表为空")
	}
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("共订阅%d个源，订阅列表\n", len(sources)))
	count := 0
	for i := range sources {
		msg.WriteString(fmt.Sprintf("[[%d]] [%s](%s)\n", sources[i].ID, sources[i].Title, sources[i].Link))
		count++
		if count == MaxSubsSizePerPage {
			ctx.Send(msg.String(), &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown})
			count = 0
			msg.Reset()
		}
	}

	if count != 0 {
		ctx.Send(msg.String(), &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown})
	}
	return nil
}
