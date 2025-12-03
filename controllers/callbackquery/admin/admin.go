package admin

import (
	"context"
	"fmt"

	"coffee-helper/controllers/middleware"
	"coffee-helper/models"
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

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "newAdminInvite", bot.MatchTypePrefix, c.newAdminInvite, c.middleware.Admin)

	return c
}

func (c *Controller) newAdminInvite(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	sentFrom := c.middleware.SentFrom(ctx)
	roleIssue, err := c.services.User.Role.GetOrCreateIssue(c.middleware.GetUser(ctx), models.FlagUserAdmin)
	if err != nil {
		panic(err)
	}
	if roleIssue == nil {
		panic("не удалось создать новый role issue. По идее такого быть не должно, но навсякий поставлю обработку на эту ошибку, а то мало ли что)\nЕсли это кто-то увидел, напишите пж @golangify, что руки кривые. Спасибо 😆")
	}
	if _, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: sentFrom.ID,
		Text:   fmt.Sprintf("Тот, кому вы хотите выдать роль админа должен отправить мне эту команду:\n\n/roleissue%s", roleIssue.Secret),
	}); err != nil {
		panic(err)
	}
}
