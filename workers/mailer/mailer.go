package mailer

import (
	"context"
	"log"
	"time"

	"coffee-helper/models"
	"coffee-helper/services"

	"github.com/go-telegram/bot"
)

type Worker struct {
	bot      *bot.Bot
	services *services.Services
}

func New(bot *bot.Bot, services *services.Services) *Worker {
	w := &Worker{
		bot:      bot,
		services: services,
	}

	return w
}

func (w *Worker) Mail(ctx context.Context, usersList []models.User, msg bot.SendMessageParams) error {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for _, user := range usersList {
		<-t.C
		msg.ChatID = user.TgID
		log.Println(user)
		if _, err := w.bot.SendMessage(ctx, &msg); err != nil {
			// тут обработчик на заблокированность бота
			log.Println("[MAILER]", err)
			continue
		}
	}
	return nil
}

func (w *Worker) Admins(ctx context.Context, msg bot.SendMessageParams) error {
	for page := 0; ; page++ {
		adminsList, err := w.services.User.AdminsList(page)
		if err != nil {
			return err
		}
		if err := w.Mail(ctx, adminsList.Users, msg); err != nil {
			return err
		}
		if adminsList.Pagination.NextPage == 0 {
			return nil
		}
	}
}
