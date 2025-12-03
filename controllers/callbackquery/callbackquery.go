package callbackquery

import (
	"context"

	"coffee-helper/controllers/callbackquery/admin"
	"coffee-helper/controllers/middleware"
	"coffee-helper/renderers"
	"coffee-helper/services"
	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

type Controller struct {
	services   *services.Services
	middleware *middleware.Controller
	renderers  *renderers.Renderers
}

func New(b *bot.Bot, services *services.Services, middleware *middleware.Controller, renderers *renderers.Renderers) *Controller {
	c := &Controller{
		services:   services,
		middleware: middleware,
		renderers:  renderers,
	}

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "deleteMessage", bot.MatchTypePrefix, c.deleteMessage)

	_ = admin.New(b, services, middleware, renderers)

	return c
}

func (c *Controller) deleteMessage(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	messageToDelete := update.CallbackQuery.Message.Message
	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    messageToDelete.Chat.ID,
		MessageID: messageToDelete.ID,
	}); err != nil {
		panic(err)
	}
}
