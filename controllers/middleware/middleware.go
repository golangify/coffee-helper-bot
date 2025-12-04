package middleware

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"coffee-helper/config"
	"coffee-helper/models"
	"coffee-helper/renderers"
	"coffee-helper/services"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

const (
	KeyCtxSentFrom = "sentFrom"
	KeyCtxUser     = "user"
)

type Controller struct {
	config    *config.Config
	services  *services.Services
	renderers *renderers.Renderers

	recentActiveUsers sync.Map
}

func New(config *config.Config, services *services.Services, renderers *renderers.Renderers) *Controller {
	c := &Controller{
		config:    config,
		services:  services,
		renderers: renderers,
	}

	return c
}

// аутентификация пользователя, валидация активности, перехват паники
func (c *Controller) Common(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		var sentFrom *botmodels.User

		// recover при ошибках
		defer func() {
			if !c.config.IsRelease {
				return
			}
			err := recover()
			if err == nil {
				return
			}
			log.Println("[PANIC]", err)
			if sentFrom != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: sentFrom.ID,
					Text:   fmt.Sprintf("[PANIC] произошла ошибка: %v", err),
				})
			}
		}()

		if update.Message != nil {
			sentFrom = update.Message.From
			// } else if update.InlineQuery != nil {
			// sentFrom = update.InlineQuery.From
		} else if update.CallbackQuery != nil {
			sentFrom = &update.CallbackQuery.From
		}

		if sentFrom == nil {
			// не обрабатываем действия, которые не являются сообщениями или нажатиями на кнопки
			// TODO: обрабатывать блокировку бота пользователями, она как раз приходит сюда
			return
		}

		{
			// ограничение на количество запросов (rate limit) за промежуток времени config.DelayBetweenActivity
			now := time.Now()
			isValid := true
			if lastActivityTime, ok := c.recentActiveUsers.Load(sentFrom.ID); ok {
				if now.Sub(lastActivityTime.(time.Time)) < c.config.DelayBetweenActivity {
					isValid = false
				}
			}
			c.recentActiveUsers.Store(sentFrom.ID, now)
			if !isValid {
				return
			}
		}

		ctx = context.WithValue(ctx, KeyCtxSentFrom, sentFrom)

		// получаем пользователя из базы данных
		user, err := c.services.User.ByTgID(sentFrom.ID)
		if err != nil {
			panic(err)
		}
		// если	пользователя нет в базе - добавляем
		if user == nil {
			user = &models.User{
				TgID:      sentFrom.ID,
				FirstName: sentFrom.FirstName,
				LastName:  sentFrom.LastName,
				Username:  sentFrom.Username,
			}
			if err = c.services.User.New(user); err != nil {
				panic(err)
			}
		}

		// если пользователь забанен, не обрабатываем его действия
		if user.Flags.Has(models.FlagUserBanned) {
			return
		}

		ctx = context.WithValue(ctx, KeyCtxUser, user)

		next(ctx, b, update)
	}
}

// проверка на админ права
func (c *Controller) Admin(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		if c.GetUser(ctx).IsAdmin() {
			next(ctx, b, update)
		} else {
			b.SendMessage(ctx, c.renderers.Bot.Error.Message.Error(c.SentFrom(ctx), "для этого действия у вас должна быть роль администратора"))
		}
	}
}

// проверка на права редактора
func (c *Controller) Editor(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		user := c.GetUser(ctx)
		if user.IsAdmin() || user.IsEditor() {
			next(ctx, b, update)
		} else {
			b.SendMessage(ctx, c.renderers.Bot.Error.Message.Error(c.SentFrom(ctx), "для этого действия у вас должна быть роль редактора"))
		}
	}
}

// проверка на права обычного сотрудника
func (c *Controller) Staff(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		user := c.GetUser(ctx)
		if user.IsAdmin() || user.IsEditor() || user.IsStaff() {
			next(ctx, b, update)
		} else {
			b.SendMessage(ctx, c.renderers.Bot.Error.Message.Error(c.SentFrom(ctx), "для этого действия у вас должна быть роль сотрудника"))
		}
	}
}

// получить пользователя из контекста
func (c *Controller) GetUser(ctx context.Context) *models.User {
	return ctx.Value(KeyCtxUser).(*models.User)
}

// получить профиль телеграм из контекста
func (c *Controller) SentFrom(ctx context.Context) *botmodels.User {
	return ctx.Value(KeyCtxSentFrom).(*botmodels.User)
}
