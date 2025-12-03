package command

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"coffee-helper/config"
	"coffee-helper/controllers/command/admin"
	"coffee-helper/controllers/middleware"
	"coffee-helper/renderers"
	"coffee-helper/services"
	"coffee-helper/workers"

	"coffee-helper/models"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

const (
	GetAccessCommand = "/getaccess"
	RoleIssueCommand = "/roleissue"
)

var (
	DelayBetweenAccessRequests = time.Hour
)

type Controller struct {
	services   *services.Services
	workers    *workers.Workers
	middleware *middleware.Controller
	renderers  *renderers.Renderers

	// TODO: переписать на sync.Map
	recentAccessRequests map[uint]time.Time
	mu                   sync.Mutex
}

func New(config *config.Config, b *bot.Bot, services *services.Services, workers *workers.Workers, middleware *middleware.Controller, renderers *renderers.Renderers) *Controller {
	c := &Controller{
		services:   services,
		workers:    workers,
		middleware: middleware,
		renderers:  renderers,

		recentAccessRequests: make(map[uint]time.Time),
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, c.start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, c.help)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/getaccess", bot.MatchTypeExact, c.getaccess)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/roleissue", bot.MatchTypePrefix, c.getaccess)

	_ = admin.New(config, b, middleware, renderers)

	return c
}

func (c *Controller) start(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	if _, err := b.SendMessage(ctx, c.renderers.Bot.Message.Start(
		c.middleware.SentFrom(ctx),
		c.middleware.GetUser(ctx),
	)); err != nil {
		panic(err)
	}
}

func (c *Controller) help(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	if _, err := b.SendMessage(ctx, c.renderers.Bot.Message.Start(
		c.middleware.SentFrom(ctx),
		c.middleware.GetUser(ctx),
	)); err != nil {
		panic(err)
	}
}

func (c *Controller) getaccess(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// вынести в middleware?
	user := c.middleware.GetUser(ctx)
	if user.Flags.Has(models.FlagUserStaff) {
		panic("у вас уже есть роль сотрудника")
	}

	now := time.Now()
	if t, ok := c.recentAccessRequests[user.ID]; ok {
		waitFor := now.Sub(t)
		if waitFor < DelayBetweenAccessRequests {
			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: user.TgID,
				Text:   fmt.Sprintf("Количество запросов на получение доступа ограничено.\n\n(%v, осталось еще %v)", DelayBetweenAccessRequests, DelayBetweenAccessRequests-waitFor),
			}); err != nil {
				panic(err)
			}
			return
		}
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: user.TgID,
		Text:   "Подождите, идёт создание запроса (занимает до минуты)...",
	}); err != nil {
		panic(err)
	}

	if err := c.workers.Mailer.Admins(ctx, bot.SendMessageParams{
		Text: fmt.Sprintf("%v (TODO) запрашивает доступ к боту.", *user),
	}); err != nil {
		panic(err)
	}
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: user.TgID,
		Text:   "Запрос на получение доступа был доставлен администраторам. Ожидайте.",
	}); err != nil {
		panic(err)
	}

	c.recentAccessRequests[user.ID] = now
}

func (c *Controller) roleissue(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	user := c.middleware.GetUser(ctx)
	secret := strings.TrimPrefix(update.Message.Text, RoleIssueCommand)
	roleIssue, err := c.services.User.Role.GrantIssue(user, secret)
	if err != nil {
		if roleIssue != nil {
			panic(fmt.Sprintf("не удалось выдать роль «%s»", models.UserFlagTitle[roleIssue.RoleFlag]))
		} else {
			panic(fmt.Sprintf("ошибка выдачи роли: %s", err))
		}
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: roleIssue.Initiator.ID,
		Text:   fmt.Sprintf("Роль «%s» успешно выдана %v(TODO отрисовать юзера)", models.UserFlagTitle[roleIssue.RoleFlag], user),
	}); err != nil {
		panic(err)
	}
}
