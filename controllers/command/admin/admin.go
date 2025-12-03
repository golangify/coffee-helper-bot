package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"coffee-helper/config"
	"coffee-helper/controllers/middleware"
	"coffee-helper/renderers"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

const SetModeCommand = "/setmode"

type Controller struct {
	config     *config.Config
	middleware *middleware.Controller
	renderers  *renderers.Renderers
}

func New(config *config.Config, b *bot.Bot, middleware *middleware.Controller, renderers *renderers.Renderers) *Controller {
	c := &Controller{
		config:     config,
		middleware: middleware,
		renderers:  renderers,
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/apanel", bot.MatchTypeExact, c.apanel, c.middleware.Admin)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/setmode", bot.MatchTypePrefix, c.setmode)

	return c
}

func (c *Controller) apanel(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	}); err != nil {
		panic(err)
	}
	apanelMessage, err := b.SendMessage(ctx, c.renderers.User.Message.Apanel(
		c.middleware.SentFrom(ctx),
		c.middleware.GetUser(ctx),
	))
	if err != nil {
		panic(err)
	}

	timer := time.NewTimer(time.Minute * 15)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-ctx.Done():
	}
	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    apanelMessage.Chat.ID,
		MessageID: apanelMessage.ID,
	}); err != nil {
		if !errors.Is(err, bot.ErrorNotFound) {
			panic(err)
		}
	}
}

func (c *Controller) setmode(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	mode := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, SetModeCommand))
	switch mode {
	case "debug":
		c.config.IsRelease = false
	case "release":
		c.config.IsRelease = true
	default:
		panic("можно указать только release или debug")
	}
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: c.middleware.SentFrom(ctx).ID,
		Text:   fmt.Sprintf("установлен режим %s", mode),
	}); err != nil {
		panic(err)
	}
}
